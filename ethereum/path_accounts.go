package ethereum

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/sethvargo/go-diceware/diceware"
)

type Account struct {
	Address      string `json:"address"` // Ethereum account address derived from the key
	Passphrase   string `json:"passphrase"`
	URL          string `json:"url"`
	RPC          string `json:"rpc_url"`
	ChainID      string `json:"chain_id"`
	JSONKeystore []byte `json:"json_keystore"`
}

func accountsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "accounts/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathAccountsList,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name"),
			HelpSynopsis: "Create an Ethereum account using a generated or provided passphrase",
			HelpDescription: `

Creates (or updates) an Ethereum externally owned account (EOAs): an account controlled by a private key. Also
creates a geth compatible keystore that is protected by a passphrase that can be supplied or optionally
generated. The generator produces a high-entropy passphrase with the provided length and requirements.
The passphrase is not returned, but it is stored at a separate path (accounts/<name>/passphrase) to allow fine
grained access controls over exposure of the passphrase. The update operation will create a new keystore using
the new passphrase.

`,
			Fields: map[string]*framework.FieldSchema{
				"rpc_url": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The RPC URL for the Ethereum node.",
					Default:     "localhost:8545",
				},
				"chain_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The Ethereum network that is being used.",
					Default:     "4", // Rinkeby
				},
				"passphrase": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The passphrase used to encrypt the private key.",
				},
				"generate_passphrase": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "Generate the passphrase.",
					Default:     false,
				},
				"words": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Number of words for the passphrase.",
					Default:     6,
				},
				"separator": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Character to separate words in passphrase.",
					Default:     "-",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.pathAccountsRead,
				logical.CreateOperation: b.pathAccountsCreate,
				logical.UpdateOperation: b.pathAccountsUpdate,
				logical.DeleteOperation: b.pathAccountsDelete,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/passphrase",
			HelpSynopsis: "Read an Ethereum account's passphrase",
			HelpDescription: `

Passphrases are use to encrypt keystores - to protect private keys. The private key is always
stored as an encrypted JSON keystore, so to use it you need to decrypt it.

`,
			Fields: map[string]*framework.FieldSchema{
				"passphrase": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The passphrase used to encrypt the private key.",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathPassphraseRead,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/sign-contract",
			HelpSynopsis: "Sign and create an Ethereum contract transaction",
			HelpDescription: `

Sign and create an Ethereum contract transaction from a given Ethereum account.

`,
			Fields: map[string]*framework.FieldSchema{
				"transaction_data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The transaction data.",
				},
				"value": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Value in ETH.",
				},
				"nonce": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The transaction nonce.",
					Default:     "1",
				},
				"gas_price": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The price in gas for the transaction.",
				},
				"gas_limit": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The gas limit for the transaction.",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathTransactionSign,
			},
		},
	}
}

func (b *backend) pathPassphraseRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathPassphraseRead", "path", req.Path)
	prunedPath := strings.Replace(req.Path, "/passphrase", "", -1)
	entry, err := req.Storage.Get(prunedPath)
	var account Account
	err = entry.DecodeJSON(&account)

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"passphrase": account.Passphrase,
		},
	}, nil
}

func (b *backend) pathAccountsRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathAccountsRead", "path", req.Path)
	entry, err := req.Storage.Get(req.Path)
	var account Account
	err = entry.DecodeJSON(&account)

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"address":  account.Address,
			"chain_id": account.ChainID,
			"keystore": fmt.Sprintf("%s", account.JSONKeystore),
			"rpc_url":  account.RPC,
		},
	}, nil
}

func (b *backend) pathAccountsCreate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathAccountsCreate", "path", req.Path)
	rpc := data.Get("rpc_url").(string)
	chainID := data.Get("chain_id").(string)
	passphrase := data.Get("passphrase").(string)
	generatePassphrase := data.Get("generate_passphrase").(bool)
	words := data.Get("words").(int)
	separator := data.Get("separator").(string)
	if generatePassphrase {
		list, _ := diceware.Generate(words)
		passphrase = strings.Join(list, separator)
	}
	tmpDir, err := b.createTemporaryKeystore(req.Path)
	if err != nil {
		return nil, err
	}
	ks := keystore.NewKeyStore(tmpDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(passphrase)
	if err != nil {
		return nil, err
	}
	keystorePath := strings.Replace(account.URL.String(), PROTOCOL_KEYSTORE, "", -1)
	jsonKeystore, _ := b.readJSONKeystore(keystorePath)
	accountJSON := &Account{Address: account.Address.Hex(),
		RPC:          rpc,
		ChainID:      chainID,
		Passphrase:   passphrase,
		URL:          account.URL.String(),
		JSONKeystore: jsonKeystore}
	entry, _ := logical.StorageEntryJSON(req.Path, accountJSON)
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}
	b.removeTemporaryKeystore(req.Path)
	return &logical.Response{
		Data: map[string]interface{}{
			"account":  accountJSON.Address,
			"chain_id": accountJSON.ChainID,
			"keystore": fmt.Sprintf("%s", jsonKeystore),
			"rpc_url":  accountJSON.RPC,
		},
	}, nil
}

func (b *backend) pathAccountsUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathAccountsUpdate", "path", req.Path)
	passphrase := data.Get("passphrase").(string)
	generatePassphrase := data.Get("generate_passphrase").(bool)
	words := data.Get("words").(int)
	separator := data.Get("separator").(string)

	if generatePassphrase {
		list, _ := diceware.Generate(words)
		passphrase = strings.Join(list, separator)
	} else if passphrase == "" {
		return nil, fmt.Errorf("nothing to update - no passphrase supplied")
	}
	entry, err := req.Storage.Get(req.Path)
	if err != nil {
		return nil, err
	}
	var account Account
	err = entry.DecodeJSON(&account)

	if err != nil {
		return nil, err
	}
	_, err = b.createTemporaryKeystore(req.Path)
	if err != nil {
		return nil, err
	}
	keystorePath := strings.Replace(account.URL, PROTOCOL_KEYSTORE, "", -1)
	b.writeTemporaryKeystoreFile(keystorePath, account.JSONKeystore)

	jsonKeystore, err := b.rekeyJSONKeystore(keystorePath, account.Passphrase, passphrase)
	b.removeTemporaryKeystore(req.Path)
	if err != nil {
		return nil, err
	} else {
		account.Passphrase = passphrase
		account.JSONKeystore = jsonKeystore
		entry, _ = logical.StorageEntryJSON(req.Path, account)

		err = req.Storage.Put(entry)
		if err != nil {
			return nil, err
		}
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"address":  account.Address,
			"chain_id": account.ChainID,
			"keystore": fmt.Sprintf("%s", account.JSONKeystore),
			"rpc_url":  account.RPC,
		},
	}, nil
}

func (b *backend) pathAccountsDelete(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathAccountsDelete", "path", req.Path)
	if err := req.Storage.Delete(req.Path); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathAccountsList(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathAccountsList", "path", req.Path)
	vals, err := req.Storage.List("accounts/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

func (b *backend) pathTransactionSign(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathTransactionSign", "path", req.Path)

	value := math.MustParseBig256(data.Get("value").(string))
	nonce := math.MustParseUint64(data.Get("nonce").(string))
	gasPrice := math.MustParseBig256(data.Get("gas_price").(string))
	gasLimit := math.MustParseBig256(data.Get("gas_limit").(string))
	input := []byte(data.Get("transaction_data").(string))

	prunedPath := strings.Replace(req.Path, "/sign-contract", "", -1)
	accountEntry, _ := req.Storage.Get(prunedPath)
	var account Account
	err := accountEntry.DecodeJSON(&account)

	if err != nil {
		return nil, err
	}
	if accountEntry == nil {
		return nil, nil
	}
	chainID := math.MustParseBig256(account.ChainID)
	_, err = b.createTemporaryKeystore(req.Path)
	if err != nil {
		return nil, err
	}
	keystorePath := strings.Replace(account.URL, PROTOCOL_KEYSTORE, "", -1)
	b.writeTemporaryKeystoreFile(keystorePath, account.JSONKeystore)
	key, _ := b.readKeyFromJSONKeystore(keystorePath, account.Passphrase)
	transactor := b.ContractTransactor(key.PrivateKey)
	var rawTx *types.Transaction
	rawTx = types.NewContractCreation(nonce, value, gasLimit, gasPrice, input)
	signedTx, err := transactor.Signer(types.NewEIP155Signer(chainID), common.HexToAddress(account.Address), rawTx)
	if err != nil {
		return nil, err
	}
	encoded, _ := rlp.EncodeToBytes(signedTx)
	hexutil.Encode(encoded[:])
	b.removeTemporaryKeystore(req.Path)
	return &logical.Response{
		Data: map[string]interface{}{
			"signed_tx": hexutil.Encode(encoded[:]),
		},
	}, nil
}

func (b *backend) ContractTransactor(key *ecdsa.PrivateKey) *bind.TransactOpts {
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}
