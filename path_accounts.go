// Copyright © 2018 Immutability, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"regexp"
	"strconv"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/shopspring/decimal"
)

const (
	// Empty is the empty string
	Empty string = ""
)

// Account is an Ethereum account
type Account struct {
	Address            string   `json:"address"` // Ethereum account address derived from the private key
	PrivateKey         string   `json:"private_key"`
	PublicKey          string   `json:"public_key"` // Ethereum public key derived from the private key
	Passphrase         string   `json:"passphrase"`
	Whitelist          []string `json:"whitelist"`
	Blacklist          []string `json:"blacklist"`
	SpendingLimitTx    string   `json:"spending_limit_tx"`
	SpendingLimitTotal string   `json:"spending_limit_total"`
	TotalSpend         string   `json:"total_spend"`
}

func accountsPaths(b *EthereumBackend) []*framework.Path {
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
				"name": &framework.FieldSchema{Type: framework.TypeString},
				"whitelist": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "The list of accounts that this account can send ETH to.",
				},
				"blacklist": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: "The list of accounts that this account can't send ETH to.",
				},
				"spending_limit_tx": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The total amount of Wei allowed to be spent in a single transaction",
					Default:     "0",
				},
				"spending_limit_total": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The total amount of Wei allowed to be spent for this account",
					Default:     "0",
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
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/debit",
			HelpSynopsis: "Send ETH from an account. ",
			HelpDescription: `

Send ETH from an account.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
				"address_to": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The address of the account to send ETH to.",
				},
				"amount": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Amount of ETH (in wei).",
				},
				"gas_limit": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The gas limit for the transaction - defaults to 21000.",
					Default:     "21000",
				},
				"gas_price": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The gas price for the transaction in wei.",
					Default:     "0",
				},
				"send": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "Send the transaction to the network.",
					Default:     true,
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathDebit,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/transfer",
			HelpSynopsis: "Transfer ERC20 tokens.",
			HelpDescription: `

Transfer ERC20 tokens.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
				"address_to": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The address of the account to send ERC20 tokens to.",
				},
				"token_address": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The address of the token contract.",
				},
				"amount": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Amount of ERC20 tokens to send.",
				},
				"gas_limit": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The gas limit for the transaction.",
				},
				"gas_price": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The gas price for the transaction in wei.",
					Default:     "0",
				},
				"send": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "Send the transaction to the network.",
					Default:     true,
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathTransfer,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/sign",
			HelpSynopsis: "Sign data",
			HelpDescription: `

Sign data using a given Ethereum account.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
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
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/sign-tx",
			HelpSynopsis: "Sign raw transaction data in json format.",
			HelpDescription: `
				Sign transaction json data using a given Ethereum account. 
				Result will be a raw transaction ready for publishing in network.`,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
				"nonce": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The nonce for the transaction.",
				},
				"address_to": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The address of the account to send ETH to.",
				},
				"value": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Value of ETH (in wei).",
				},
				"tx_data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Transaction data in HEX string.",
					Default:     "0x",
				},
				"gas_limit": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The gas limit for the transaction - defaults to 21000.",
					Default:     "21000",
				},
				"gas_price": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The gas price for the transaction in wei.",
					Default:     "0",
				},
				"chain_id": &framework.FieldSchema{
					Type: framework.TypeString,
					Description: `Ethereum network - can be one of the following values:

					1 - Ethereum mainnet
					2 - Morden (disused), Expanse mainnet
					3 - Ropsten
					4 - Rinkeby
					30 - Rootstock mainnet
					31 - Rootstock testnet
					42 - Kovan
					61 - Ethereum Classic mainnet
					62 - Ethereum Classic testnet`,
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSignTx,
			},
		},
	}
}

func (b *EthereumBackend) pathAccountsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vals, err := req.Storage.List(ctx, "accounts/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

func (b *EthereumBackend) pathAccountsDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("Error reading account")
	}
	if account == nil {
		return nil, nil
	}
	if err := req.Storage.Delete(ctx, req.Path); err != nil {
		return nil, err
	}
	b.removeCrossReference(ctx, req, name, account.Address)
	return nil, nil
}

func (b *EthereumBackend) readAccount(ctx context.Context, req *logical.Request, name string) (*Account, error) {
	path := fmt.Sprintf("accounts/%s", name)
	entry, err := req.Storage.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var account Account
	err = entry.DecodeJSON(&account)

	if entry == nil {
		return nil, fmt.Errorf("failed to deserialize account at %s", path)
	}

	return &account, nil
}

func (b *EthereumBackend) pathAccountsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("Error reading account")
	}
	if account == nil {
		return nil, nil
	}
	balance, exchangeValue, err := b.readAccountBalanceByAddress(ctx, req, account.Address)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"address":              account.Address,
			"whitelist":            account.Whitelist,
			"blacklist":            account.Blacklist,
			"spending_limit_tx":    account.SpendingLimitTx,
			"spending_limit_total": account.SpendingLimitTotal,
			"total_spend":          account.TotalSpend,
			"balance":              balance,
			"balance_in_usd":       exchangeValue,
		},
	}, nil
}

func (b *EthereumBackend) pathAccountsCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	spendingLimitTxString := data.Get("spending_limit_tx").(string)
	spendingLimitTx, err := decimal.NewFromString(spendingLimitTxString)
	if err != nil || spendingLimitTx.IsNegative() {
		return nil, fmt.Errorf("spending_limit_tx is either not a number or is negative")
	}
	spendingLimitTotalString := data.Get("spending_limit_total").(string)
	spendingLimitTotal, err := decimal.NewFromString(spendingLimitTotalString)
	if err != nil || spendingLimitTotal.IsNegative() {
		return nil, fmt.Errorf("spending_limit_tx is either not a number or is negative")
	}
	var whiteList []string
	if whiteListRaw, ok := data.GetOk("whitelist"); ok {
		whiteList = whiteListRaw.([]string)
	}
	var blackList []string
	if blackListRaw, ok := data.GetOk("blacklist"); ok {
		blackList = blackListRaw.([]string)
	}

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	defer ZeroKey(privateKey)
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyString := hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyString := hexutil.Encode(publicKeyBytes)[4:]

	hash := sha3.NewKeccak256()
	hash.Write(publicKeyBytes[1:])
	address := hexutil.Encode(hash.Sum(nil)[12:])

	accountJSON := &Account{
		Address:            address,
		PrivateKey:         privateKeyString,
		PublicKey:          publicKeyString,
		Whitelist:          Dedup(whiteList),
		Blacklist:          Dedup(blackList),
		SpendingLimitTx:    spendingLimitTx.String(),
		SpendingLimitTotal: spendingLimitTotal.String(),
		TotalSpend:         "0",
	}
	entry, err := logical.StorageEntryJSON(req.Path, accountJSON)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}
	b.crossReference(ctx, req, name, accountJSON.Address)
	return &logical.Response{
		Data: map[string]interface{}{
			"address":              accountJSON.Address,
			"whitelist":            accountJSON.Whitelist,
			"blacklist":            accountJSON.Blacklist,
			"spending_limit_tx":    accountJSON.SpendingLimitTx,
			"spending_limit_total": accountJSON.SpendingLimitTotal,
			"total_spend":          accountJSON.TotalSpend,
		},
	}, nil
}

func (b *EthereumBackend) pathAccountUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("Error reading account")
	}
	if account == nil {
		return nil, nil
	}

	spendingLimitTxString := data.Get("spending_limit_tx").(string)
	spendingLimitTx, err := decimal.NewFromString(spendingLimitTxString)
	if err != nil || spendingLimitTx.IsNegative() {
		return nil, fmt.Errorf("spending_limit_tx is either not a number or is negative")
	}
	spendingLimitTotalString := data.Get("spending_limit_total").(string)
	spendingLimitTotal, err := decimal.NewFromString(spendingLimitTotalString)
	if err != nil || spendingLimitTotal.IsNegative() {
		return nil, fmt.Errorf("spending_limit_tx is either not a number or is negative")
	}
	var whiteList []string
	if whiteListRaw, ok := data.GetOk("whitelist"); ok {
		whiteList = whiteListRaw.([]string)
	}
	var blackList []string
	if blackListRaw, ok := data.GetOk("blacklist"); ok {
		blackList = blackListRaw.([]string)
	}

	// Update the account
	account.SpendingLimitTx = spendingLimitTx.String()
	account.SpendingLimitTotal = spendingLimitTotal.String()
	account.Whitelist = whiteList
	account.Blacklist = blackList

	entry, err := logical.StorageEntryJSON(req.Path, account)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"address":              account.Address,
			"whitelist":            account.Whitelist,
			"blacklist":            account.Blacklist,
			"spending_limit_tx":    account.SpendingLimitTx,
			"spending_limit_total": account.SpendingLimitTotal,
			"total_spend":          account.TotalSpend,
		},
	}, nil

}

func (b *EthereumBackend) validAccountConstraints(account *Account, amount *big.Int, toAddress string) (bool, error) {
	txLimit := ValidNumber(account.SpendingLimitTx)
	limit := ValidNumber(account.SpendingLimitTotal)
	totalSpend := ValidNumber(account.TotalSpend)

	if txLimit.Cmp(amount) == -1 && txLimit.Cmp(big.NewInt(0)) == 1 {
		return false, fmt.Errorf("transaction amount (%s) is larger than the transactional limit (%s)", amount.String(), account.SpendingLimitTx)
	}

	if limit.Cmp(totalSpend.Add(totalSpend, amount)) == -1 && limit.Cmp(big.NewInt(0)) == 1 {
		return false, fmt.Errorf("transaction amount (%s) + total spend (%s) is larger than the limit (%s)", amount.String(), account.TotalSpend, account.SpendingLimitTotal)
	}

	if contains(account.Blacklist, toAddress) {
		return false, fmt.Errorf("%s is blacklisted", toAddress)
	}

	if len(account.Whitelist) > 0 && !contains(account.Whitelist, toAddress) {
		return false, fmt.Errorf("%s is not in the whitelist", toAddress)
	}

	return true, nil
}

func (b *EthereumBackend) updateTotalSpend(ctx context.Context, req *logical.Request, path string, account *Account, amount *big.Int) (string, error) {
	startingPoint := account.TotalSpend
	totalSpend := ValidNumber(account.TotalSpend)
	totalSpend = totalSpend.Add(totalSpend, amount)
	account.TotalSpend = totalSpend.String()
	entry, err := logical.StorageEntryJSON(path, account)
	if err != nil {
		return startingPoint, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return startingPoint, err
	}
	return account.TotalSpend, nil
}

func (b *EthereumBackend) verifySignature(ctx context.Context, req *logical.Request, data *framework.FieldData, name string) (*logical.Response, error) {
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("error reading account")
	}
	if account == nil {
		return nil, nil
	}
	signature := data.Get("signature").(string)
	dataToSign := data.Get("data").(string)
	privateKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer ZeroKey(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	dataBytes := []byte(dataToSign)
	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(dataBytes)

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signatureBytes)
	if err != nil {
		return nil, err
	}

	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	if !matches {
		return nil, fmt.Errorf("signature not verified")
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"verified":  matches,
			"signature": signature,
			"address":   account.Address,
		},
	}, nil

}

func (b *EthereumBackend) pathVerify(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)
	return b.verifySignature(ctx, req, data, name)
}

// ValidNumber returns a valid positive integer
func ValidNumber(input string) *big.Int {
	if input == "" {
		return big.NewInt(0)
	}
	matched, err := regexp.MatchString("([0-9])", input)
	if !matched || err != nil {
		return nil
	}
	amount := math.MustParseBig256(input)
	return amount.Abs(amount)
}

func (b *EthereumBackend) pathDebit(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)
	sendTransaction := data.Get("send").(bool)
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("error reading account")
	}
	if account == nil {
		return nil, nil
	}
	balance, _, exchangeValue, err := b.readAccountBalance(ctx, req, name)
	if err != nil {
		return nil, err
	}
	amount := ValidNumber(data.Get("amount").(string))
	if amount == nil {
		return nil, fmt.Errorf("invalid amount")
	}
	if amount.Cmp(balance) > 0 {
		return nil, fmt.Errorf("Insufficient funds spend %v because the current account balance is %v", amount, balance)
	}
	if valid, err := b.validAccountConstraints(account, amount, data.Get("address_to").(string)); !valid {
		return nil, err
	}
	chainID := ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}
	gasLimitIn := ValidNumber(data.Get("gas_limit").(string))
	if gasLimitIn == nil {
		return nil, fmt.Errorf("invalid gas limit")
	}
	gasLimit := gasLimitIn.Uint64()

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to " + config.getRPCURL())
	}
	gasPrice := ValidNumber(data.Get("gas_price").(string))
	if big.NewInt(0).Cmp(gasPrice) == 0 {
		gasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, err
		}
	}
	privateKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error reconstructing private key")
	}
	defer ZeroKey(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	toAddress := common.HexToAddress(data.Get("address_to").(string))
	var txData []byte
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, txData)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}
	if sendTransaction {
		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return nil, err
		}
	}

	totalSpend, err := b.updateTotalSpend(ctx, req, fmt.Sprintf("accounts/%s", name), account, amount)
	if err != nil {
		return nil, err
	}
	amountInUSD, _ := decimal.NewFromString("0")
	if config.ChainID == EthereumMainnet {
		amountInUSD, err = ConvertToUSD(amount.String())
		if err != nil {
			return nil, err
		}
	}
	var signedTxBuff bytes.Buffer
	signedTx.EncodeRLP(&signedTxBuff)

	return &logical.Response{
		Data: map[string]interface{}{
			"transaction_hash":        signedTx.Hash().Hex(),
			"signed_transaction":      hexutil.Encode(signedTxBuff.Bytes()),
			"address_from":            account.Address,
			"address_to":              toAddress.String(),
			"amount":                  amount.String(),
			"amount_in_usd":           amountInUSD,
			"gas_price":               gasPrice.String(),
			"gas_limit":               gasLimitIn.String(),
			"total_spend":             totalSpend,
			"starting_balance":        balance,
			"starting_balance_in_usd": exchangeValue,
		},
	}, nil
}

func (b *EthereumBackend) pathSign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)
	dataToSign := data.Get("data").(string)
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("error reading account")
	}
	if account == nil {
		return nil, nil
	}

	privateKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer ZeroKey(privateKey)
	dataBytes := []byte(dataToSign)
	hash := crypto.Keccak256Hash(dataBytes)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": hexutil.Encode(signature),
			"address":   account.Address,
		},
	}, nil
}

func (b *EthereumBackend) pathSignTx(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("error reading account")
	}
	if account == nil {
		return nil, nil
	}

	// Read tx parameters from json body
	value := ValidNumber(data.Get("value").(string))
	if value == nil {
		return nil, fmt.Errorf("invalid value")
	}

	chainID := ValidNumber(data.Get("chain_id").(string))
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	nonceIn := ValidNumber(data.Get("nonce").(string))
	if nonceIn == nil {
		return nil, fmt.Errorf("invalid nonce")
	}
	nonce := nonceIn.Uint64()

	gasLimitIn := ValidNumber(data.Get("gas_limit").(string))
	if gasLimitIn == nil {
		return nil, fmt.Errorf("invalid gas limit")
	}
	gasLimit := gasLimitIn.Uint64()

	gasPrice := ValidNumber(data.Get("gas_price").(string))
	if gasPrice == nil {
		return nil, fmt.Errorf("invalid gas price")
	}

	hexAddress := data.Get("address_to").(string)
	if !common.IsHexAddress(hexAddress) {
		return nil, fmt.Errorf("invalid transaction receiver address")
	}
	addressTo := common.HexToAddress(data.Get("address_to").(string))

	txData, err := hexutil.Decode(data.Get("tx_data").(string))
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer ZeroKey(privateKey)

	// Create transaction
	tx := types.NewTransaction(nonce, addressTo, value, gasLimit, gasPrice, txData)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	// Encode raw transaction and prepare for pushing in network
	var signedTxBuff bytes.Buffer
	signedTx.EncodeRLP(&signedTxBuff)

	return &logical.Response{
		Data: map[string]interface{}{
			"signed_tx": hexutil.Encode(signedTxBuff.Bytes()),
			"address":   account.Address,
		},
	}, nil
}

func (b *EthereumBackend) pathContractsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vals, err := req.Storage.List(ctx, req.Path)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

func (b *EthereumBackend) pathTransfer(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)
	sendTransaction := data.Get("send").(bool)
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("error reading account")
	}
	if account == nil {
		return nil, nil
	}

	amount := ValidNumber(data.Get("amount").(string))
	if amount == nil {
		return nil, fmt.Errorf("invalid amount")
	}
	chainID := ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to " + config.getRPCURL())
	}
	gasPrice := ValidNumber(data.Get("gas_price").(string))
	if big.NewInt(0).Cmp(gasPrice) == 0 {
		gasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, err
		}
	}
	privateKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error reconstructing private key")
	}
	defer ZeroKey(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	toAddress := common.HexToAddress(data.Get("address_to").(string))
	tokenAddress := common.HexToAddress(data.Get("token_address").(string))
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	var txData []byte
	txData = append(txData, methodID...)
	txData = append(txData, paddedAddress...)
	txData = append(txData, paddedAmount...)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: txData,
	})
	if err != nil {
		return nil, err
	}
	value := big.NewInt(0) // in wei (0 eth)

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, txData)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}
	if sendTransaction {
		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return nil, err
		}
	}
	var signedTxBuff bytes.Buffer
	signedTx.EncodeRLP(&signedTxBuff)

	return &logical.Response{
		Data: map[string]interface{}{
			"transaction_hash":   signedTx.Hash().Hex(),
			"signed_transaction": hexutil.Encode(signedTxBuff.Bytes()),
			"address_from":       account.Address,
			"address_to":         toAddress.String(),
			"token_address":      tokenAddress.String(),
			"amount":             amount.String(),
			"gas_price":          gasPrice.String(),
			"gas_limit":          strconv.FormatUint(gasLimit, 10),
		},
	}, nil
}

func (b *EthereumBackend) readAccountBalanceByAddress(ctx context.Context, req *logical.Request, address string) (*big.Int, decimal.Decimal, error) {
	zero, _ := decimal.NewFromString("0")
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, zero, err
	}
	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, zero, fmt.Errorf("cannot connect to " + config.getRPCURL())
	}
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return nil, zero, err
	}
	// Calculate exchange rate value if on Mainnet
	if config.ChainID == EthereumMainnet {
		exchangeValue, err := ConvertToUSD(balance.String())
		if err != nil {
			return nil, zero, err
		}
		return balance, exchangeValue, nil
	}
	return balance, zero, nil
}

func (b *EthereumBackend) readAccountBalance(ctx context.Context, req *logical.Request, name string) (*big.Int, string, decimal.Decimal, error) {
	zero, _ := decimal.NewFromString("0")
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, Empty, zero, fmt.Errorf("error reading account")
	}
	if account == nil {
		return nil, Empty, zero, nil
	}

	balance, exchangeValue, err := b.readAccountBalanceByAddress(ctx, req, account.Address)
	if err != nil {
		return nil, Empty, zero, err
	}
	return balance, account.Address, exchangeValue, nil
}
