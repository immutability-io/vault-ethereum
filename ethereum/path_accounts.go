package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	rpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/sethvargo/go-diceware/diceware"
)

type Account struct {
	Address        string   `json:"address"` // Ethereum account address derived from the key
	Passphrase     string   `json:"passphrase"`
	KeystoreName   string   `json:"keystore_name"`
	RPC            string   `json:"rpc_url"`
	ChainID        string   `json:"chain_id"`
	Whitelist      []string `json:"whitelist"`
	Blacklist      []string `json:"blacklist"`
	JSONKeystore   []byte   `json:"json_keystore"`
	PendingBalance *big.Int `json:"pending_balance"`
	PendingNonce   uint64   `json:"pending_nonce"`
	PendingTxCount uint     `json:"pending_tx_count"`
}

func accountsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "accounts/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathAccountsList,
			},
			HelpSynopsis: "List all the Ethereum accounts at a path",
			HelpDescription: `
			All the Ethereum accounts will be listed.
			`,
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
					Default:     "http://localhost:8545",
				},
				"chain_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The Ethereum network that is being used.",
					Default:     "4", // Rinkeby
				},
				"whitelist": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "The list of accounts that this account can send ETH to.",
				},
				"blacklist": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "The list of accounts that this account can't send ETH to.",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.pathAccountsRead,
				logical.CreateOperation: b.pathAccountsCreate,
				logical.UpdateOperation: b.pathAccountUpdate,
				logical.DeleteOperation: b.pathAccountsDelete,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/export",
			HelpSynopsis: "Export a single Ethereum JSON keystore from vault into the provided path.",
			HelpDescription: `

Writes a JSON keystore to a folder (e.g., /Users/immutability/.ethereum/keystore).

`,
			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Directory to export the keystore into - must be an absolute path.",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathExportCreate,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/balance",
			HelpSynopsis: "Convenience method to read the balance of an account.",
			HelpDescription: `

Queries the Ethereum blockchain (chain_id) for the balance of an account.

`,
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathAccountBalanceRead,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/debit",
			HelpSynopsis: "Sign and create an Ethereum contract transaction",
			HelpDescription: `

		Send ether from a given Ethereum account.

		`,
			Fields: map[string]*framework.FieldSchema{
				"to": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The address of the account to send ETH to.",
				},
				"amount": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Amount of ETH (in Gwei).",
				},
				"gas_limit": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The gas limit for the transaction.",
					Default:     "50000",
				},
				"gas_price": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The gas price for the transaction in wei.",
					Default:     "20000000000",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathDebit,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/sign",
			HelpSynopsis: "Hash and sign data",
			HelpDescription: `

Hash and sign data using a given Ethereum account.

`,
			Fields: map[string]*framework.FieldSchema{
				"data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The data to hash (keccak) and sign.",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSign,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/verify",
			HelpSynopsis: "Verify that this account signed something.",
			HelpDescription: `

Validate that this account signed some data.

`,
			Fields: map[string]*framework.FieldSchema{
				"data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The data to verify the signature of.",
				},
				"signature": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The signature to verify",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathVerify,
			},
		},
	}
}

func (b *backend) pathAccountsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	account, err := b.readAccount(ctx, req, req.Path, false)
	if err != nil {
		return nil, err
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"address":   account.Address,
			"chain_id":  account.ChainID,
			"whitelist": account.Whitelist,
			"blacklist": account.Blacklist,
			"rpc_url":   account.RPC,
		},
	}, nil
}

func (b *backend) pathAccountBalanceRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	prunedPath := strings.Replace(req.Path, "/balance", "", -1)
	account, err := b.readAccount(ctx, req, prunedPath, true)
	if err != nil {
		return nil, err
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"address":          account.Address,
			"pending_balance":  account.PendingBalance.String(),
			"pending_nonce":    fmt.Sprintf("%d", account.PendingNonce),
			"pending_tx_count": fmt.Sprintf("%d", account.PendingTxCount),
		},
	}, nil
}

func (b *backend) pathAccountsCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	rpc := data.Get("rpc_url").(string)
	chainID := data.Get("chain_id").(string)
	whitelist := data.Get("whitelist").([]string)
	blacklist := data.Get("blacklist").([]string)
	list, _ := diceware.Generate(PassphraseWords)
	passphrase := strings.Join(list, PassphraseSeparator)
	tmpDir, err := b.createTemporaryKeystoreDirectory()
	if err != nil {
		return nil, err
	}
	ks := keystore.NewKeyStore(tmpDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(passphrase)
	if err != nil {
		return nil, err
	}
	keystorePath := strings.Replace(account.URL.String(), ProtocolKeystore, "", -1)

	jsonKeystore, err := b.readJSONKeystore(keystorePath)
	if err != nil {
		return nil, err
	}
	accountJSON := &Account{Address: account.Address.Hex(),
		RPC:          rpc,
		ChainID:      chainID,
		Passphrase:   passphrase,
		Whitelist:    dedup(whitelist),
		Blacklist:    dedup(blacklist),
		KeystoreName: filepath.Base(account.URL.String()),
		JSONKeystore: jsonKeystore}
	entry, err := logical.StorageEntryJSON(req.Path, accountJSON)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}
	b.removeTemporaryKeystore(tmpDir)
	return &logical.Response{
		Data: map[string]interface{}{
			"address":   accountJSON.Address,
			"chain_id":  accountJSON.ChainID,
			"whitelist": accountJSON.Whitelist,
			"blacklist": accountJSON.Blacklist,
			"rpc_url":   accountJSON.RPC,
		},
	}, nil
}

func (b *backend) pathAccountUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	whitelist := data.Get("whitelist").([]string)
	blacklist := data.Get("blacklist").([]string)
	account, err := b.readAccount(ctx, req, req.Path, false)
	account.Whitelist = dedup(whitelist)
	account.Blacklist = dedup(blacklist)

	entry, _ := logical.StorageEntryJSON(req.Path, account)

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"address":   account.Address,
			"chain_id":  account.ChainID,
			"whitelist": account.Whitelist,
			"blacklist": account.Blacklist,
			"rpc_url":   account.RPC,
		},
	}, nil
}

func (b *backend) pathAccountsDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, req.Path); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathAccountsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vals, err := req.Storage.List(ctx, "accounts/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

func (b *backend) pathSign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	input := []byte(data.Get("data").(string))
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(input), input)
	hash := crypto.Keccak256([]byte(msg))
	prunedPath := strings.Replace(req.Path, "/sign", "", -1)
	account, err := b.readAccount(ctx, req, prunedPath, false)
	if err != nil {
		return nil, err
	}
	key, err := b.getAccountPrivateKey(prunedPath, *account)
	if err != nil {
		return nil, err
	}
	defer zeroKey(key.PrivateKey)
	signature, err := crypto.Sign(hash, key.PrivateKey)

	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"signature": hexutil.Encode(signature[:]),
		},
	}, nil
}

func (b *backend) pathVerify(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	input := []byte(data.Get("data").(string))
	signatureRaw := data.Get("signature").(string)

	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(input), input)
	hash := crypto.Keccak256([]byte(msg))
	prunedPath := strings.Replace(req.Path, "/verify", "", -1)
	account, err := b.readAccount(ctx, req, prunedPath, false)
	if err != nil {
		return nil, err
	}
	signature, err := hexutil.Decode(signatureRaw)

	if err != nil {
		return nil, err
	}
	pubkey, err := crypto.SigToPub(hash, signature)

	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(*pubkey)

	verified := account.Address == address.Hex()
	return &logical.Response{
		Data: map[string]interface{}{
			"verified": verified,
		},
	}, nil
}

func (b *backend) pathDebit(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	prunedPath := strings.Replace(req.Path, "/debit", "", -1)
	amount := math.MustParseBig256(data.Get("amount").(string))
	gasLimit := math.MustParseBig256(data.Get("gas_limit").(string))
	gasPrice := math.MustParseBig256(data.Get("gas_price").(string))
	toAddress := common.HexToAddress(data.Get("to").(string))

	account, err := b.readAccount(ctx, req, prunedPath, true)
	if err != nil {
		return nil, err
	}
	allowed, err := b.isDebitAllowed(account, data.Get("to").(string), amount)
	if !allowed {
		return nil, err
	}
	chainID := math.MustParseBig256(account.ChainID)
	key, err := b.getAccountPrivateKey(prunedPath, *account)
	if err != nil {
		return nil, err
	}
	defer zeroKey(key.PrivateKey)

	transactor := b.NewTransactor(key.PrivateKey)
	var rawTx *types.Transaction

	client, err := rpc.Dial(account.RPC)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(client)
	fromAddress := common.HexToAddress(account.Address)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	nonce, err := ethClient.NonceAt(ctx, fromAddress, nil)
	if err != nil {
		return nil, err
	}

	if !gasLimit.IsUint64() {
		return nil, errors.New("Cannot convert gas limit to uint64")
	}
	gl := gasLimit.Uint64()
	rawTx = types.NewTransaction(nonce, toAddress, amount, gl, gasPrice, nil)
	signedTx, err := transactor.Signer(types.NewEIP155Signer(chainID), common.HexToAddress(account.Address), rawTx)
	if err != nil {
		return nil, err
	}
	ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	txHash := rawTx.Hash().Hex()

	return &logical.Response{
		Data: map[string]interface{}{
			"from_address": account.Address,
			"to_address":   toAddress.String(),
			"tx_hash":      txHash,
		},
	}, nil
}

func (b *backend) pathExportCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	directory := data.Get("path").(string)
	prunedPath := strings.Replace(req.Path, "/export", "", -1)
	account, err := b.readAccount(ctx, req, prunedPath, false)
	if err != nil {
		return nil, err
	}
	list, _ := diceware.Generate(PassphraseWords)
	passphrase := strings.Join(list, PassphraseSeparator)
	tmpDir, err := b.createTemporaryKeystoreDirectory()
	if err != nil {
		return nil, err
	}
	keystorePath, err := b.writeTemporaryKeystoreFile(tmpDir, account.KeystoreName, account.JSONKeystore)
	if err != nil {
		return nil, err
	}

	jsonKeystore, err := b.rekeyJSONKeystore(keystorePath, account.Passphrase, passphrase)
	b.removeTemporaryKeystore(tmpDir)
	if err != nil {
		return nil, err
	} else {
		account.Passphrase = passphrase
		account.JSONKeystore = jsonKeystore
	}
	filePath, err := b.exportKeystore(directory, account)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"path":       filePath,
			"passphrase": passphrase,
		},
	}, nil
}
