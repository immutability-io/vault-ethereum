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
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	bip44 "github.com/immutability-io/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/immutability-io/vault-ethereum/util"
)

const (
	// DerivationPath is the root in a BIP44 hdwallet
	DerivationPath string = "m/44'/60'/0'/0/%d"
	// Empty is the empty string
	Empty string = ""
	// Utf8Encoding is utf
	Utf8Encoding string = "utf8"
	// HexEncoding is hex
	HexEncoding string = "hex"
)

// AccountJSON is what we store for an Ethereum account
type AccountJSON struct {
	Index      int      `json:"index"`
	Mnemonic   string   `json:"mnemonic"`
	Inclusions []string `json:"inclusions"`
	Exclusions []string `json:"exclusions"`
}

// ValidAddress returns an error if the address is not included or if it is excluded
func (account *AccountJSON) ValidAddress(toAddress *common.Address) error {
	if util.Contains(account.Exclusions, toAddress.Hex()) {
		return fmt.Errorf("%s is excludeded by this account", toAddress.Hex())
	}

	if len(account.Inclusions) > 0 && !util.Contains(account.Inclusions, toAddress.Hex()) {
		return fmt.Errorf("%s is not in the set of inclusions of this account", toAddress.Hex())
	}
	return nil
}

// TransactionParams are typical parameters for a transaction
type TransactionParams struct {
	Nonce    uint64          `json:"nonce"`
	Address  *common.Address `json:"address"`
	Amount   *big.Int        `json:"amount"`
	GasPrice *big.Int        `json:"gas_price"`
	GasLimit uint64          `json:"gas_limit"`
}

func accountPaths(b *PluginBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: QualifiedPath("accounts/?"),
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathAccountsList,
			},
			HelpSynopsis: "List all the Ethereum accounts at a path",
			HelpDescription: `
			All the Ethereum accounts will be listed.
			`,
		},
		{
			Pattern:      QualifiedPath("accounts/" + framework.GenericNameRegex("name")),
			HelpSynopsis: "Create an Ethereum account using a generated or provided passphrase.",
			HelpDescription: `

Creates (or updates) an Ethereum account: an account controlled by a private key. Also
The generator produces a high-entropy passphrase with the provided length and requirements.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"mnemonic": {
					Type:        framework.TypeString,
					Default:     Empty,
					Description: "The mnemonic to use to create the account. If not provided, one is generated.",
				},
				"index": {
					Type:        framework.TypeInt,
					Description: "The index used in BIP-44.",
					Default:     0,
				},
				"inclusions": {
					Type:        framework.TypeCommaStringSlice,
					Description: "The list of accounts that this account can send transactions to.",
				},
				"exclusions": {
					Type:        framework.TypeCommaStringSlice,
					Description: "The list of accounts that this account can't send transactions to.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.pathAccountsRead,
				logical.CreateOperation: b.pathAccountsCreate,
				logical.UpdateOperation: b.pathAccountUpdate,
				logical.DeleteOperation: b.pathAccountsDelete,
			},
		},
		{
			Pattern:      QualifiedPath("accounts/" + framework.GenericNameRegex("name") + "/transfer"),
			HelpSynopsis: "Send ETH from an account.",
			HelpDescription: `

Send ETH from an account.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"to": {
					Type:        framework.TypeString,
					Description: "The address of the wallet to send ETH to.",
				},
				"amount": {
					Type:        framework.TypeString,
					Description: "Amount of ETH (in wei).",
				},
				"gas_limit": {
					Type:        framework.TypeString,
					Description: "The gas limit for the transaction - defaults to 21000.",
					Default:     "21000",
				},
				"gas_price": {
					Type:        framework.TypeString,
					Description: "The gas price for the transaction in wei.",
					Default:     "0",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathTransfer,
				logical.CreateOperation: b.pathTransfer,
			},
		},
		{
			Pattern:      QualifiedPath("accounts/" + framework.GenericNameRegex("name") + "/balance"),
			HelpSynopsis: "Return the balance for an account.",
			HelpDescription: `

Return the balance in wei for an address.

`,
			Fields: map[string]*framework.FieldSchema{
				"name":    {Type: framework.TypeString},
				"address": {Type: framework.TypeString},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathReadBalance,
			},
		},
		{
			Pattern:      QualifiedPath("accounts/" + framework.GenericNameRegex("name") + "/sign-tx"),
			HelpSynopsis: "Sign a transaction.",
			HelpDescription: `

Sign a transaction.

`,
			Fields: map[string]*framework.FieldSchema{
				"name":    {Type: framework.TypeString},
				"address": {Type: framework.TypeString},
				"to": {
					Type:        framework.TypeString,
					Description: "The address of the wallet to send ETH to.",
				},
				"data": {
					Type:        framework.TypeString,
					Description: "The data to sign.",
				},
				"encoding": {
					Type:        framework.TypeString,
					Default:     "utf8",
					Description: "The encoding of the data to sign.",
				},
				"amount": {
					Type:        framework.TypeString,
					Description: "Amount of ETH (in wei).",
				},
				"nonce": {
					Type:        framework.TypeString,
					Description: "The transaction nonce.",
				},
				"gas_limit": {
					Type:        framework.TypeString,
					Description: "The gas limit for the transaction - defaults to 21000.",
					Default:     "21000",
				},
				"gas_price": {
					Type:        framework.TypeString,
					Description: "The gas price for the transaction in wei.",
					Default:     "0",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSignTx,
				logical.UpdateOperation: b.pathSignTx,
			},
		},
		{
			Pattern:      QualifiedPath("accounts/" + framework.GenericNameRegex("name") + "/deploy"),
			HelpSynopsis: "Deploy a smart contract from an account.",
			HelpDescription: `

Deploy a smart contract to the network.

`,
			Fields: map[string]*framework.FieldSchema{
				"name":    {Type: framework.TypeString},
				"address": {Type: framework.TypeString},
				"version": {
					Type:        framework.TypeString,
					Description: "The smart contract version.",
				},
				"abi": {
					Type:        framework.TypeString,
					Description: "The contract ABI.",
				},
				"bin": {
					Type:        framework.TypeString,
					Description: "The compiled smart contract.",
				},
				"gas_limit": {
					Type:        framework.TypeString,
					Description: "The gas limit for the transaction - defaults to 0 meaning estimate.",
					Default:     "0",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathDeploy,
				logical.CreateOperation: b.pathDeploy,
			},
		},
		{
			Pattern:      QualifiedPath("accounts/" + framework.GenericNameRegex("name") + "/sign"),
			HelpSynopsis: "Sign a message",
			HelpDescription: `

Sign calculates an ECDSA signature for:
keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).

https://eth.wiki/json-rpc/API#eth_sign

		`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"message": {
					Type:        framework.TypeString,
					Description: "Message to sign.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSignMessage,
				logical.UpdateOperation: b.pathSignMessage,
			},
		},
	}
}

func (b *PluginBackend) pathAccountsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vals, err := req.Storage.List(ctx, QualifiedPath("accounts/"))
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

func readAccount(ctx context.Context, req *logical.Request, name string) (*AccountJSON, error) {
	path := QualifiedPath(fmt.Sprintf("accounts/%s", name))
	entry, err := req.Storage.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var accountJSON AccountJSON
	err = entry.DecodeJSON(&accountJSON)

	if entry == nil {
		return nil, fmt.Errorf("failed to deserialize account at %s", path)
	}
	return &accountJSON, nil
}

func (b *PluginBackend) pathAccountsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	accountJSON, err := readAccount(ctx, req, name)

	_, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}
	if err != nil || accountJSON == nil {
		return nil, fmt.Errorf("Error reading account")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address":    account.Address.Hex(),
			"inclusions": accountJSON.Inclusions,
			"exclusions": accountJSON.Exclusions,
		},
	}, nil
}

func (b *PluginBackend) pathAccountsDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	_, err = readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Delete(ctx, req.Path); err != nil {
		return nil, err
	}
	return nil, nil
}

func getWalletAndAccount(accountJSON AccountJSON) (*bip44.Wallet, *accounts.Account, error) {
	hdwallet, err := bip44.NewFromMnemonic(accountJSON.Mnemonic)
	if err != nil {
		return nil, nil, err
	}
	derivationPath := fmt.Sprintf(DerivationPath, accountJSON.Index)
	path := bip44.MustParseDerivationPath(derivationPath)
	account, err := hdwallet.Derive(path, true)
	if err != nil {
		return nil, nil, err
	}
	return hdwallet, &account, nil
}

func (b *PluginBackend) pathAccountsCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	var inclusions []string
	if inclusionsRaw, ok := data.GetOk("inclusions"); ok {
		inclusions = inclusionsRaw.([]string)
	}
	var exclusions []string
	if exclusionsRaw, ok := data.GetOk("exclusions"); ok {
		exclusions = exclusionsRaw.([]string)
	}
	index := data.Get("index").(int)
	mnemonic := data.Get("mnemonic").(string)
	if mnemonic == Empty {
		entropy, err := bip39.NewEntropy(128)
		if err != nil {
			return nil, err
		}

		mnemonic, err = bip39.NewMnemonic(entropy)

	}

	if err != nil {
		return nil, err
	}
	accountJSON := &AccountJSON{
		Index:      index,
		Mnemonic:   mnemonic,
		Inclusions: util.Dedup(inclusions),
		Exclusions: util.Dedup(exclusions),
	}
	_, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	err = b.updateAccount(ctx, req, name, accountJSON)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address":    account.Address.Hex(),
			"inclusions": accountJSON.Inclusions,
			"exclusions": accountJSON.Exclusions,
		},
	}, nil
}

func (b *PluginBackend) updateAccount(ctx context.Context, req *logical.Request, name string, accountJSON *AccountJSON) error {
	path := QualifiedPath(fmt.Sprintf("accounts/%s", name))

	entry, err := logical.StorageEntryJSON(path, accountJSON)
	if err != nil {
		return err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return err
	}
	return nil
}

func (b *PluginBackend) pathAccountUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)
	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	var inclusions []string
	if inclusionsRaw, ok := data.GetOk("inclusions"); ok {
		inclusions = inclusionsRaw.([]string)
	}
	var exclusions []string
	if exclusionsRaw, ok := data.GetOk("exclusions"); ok {
		exclusions = exclusionsRaw.([]string)
	}
	accountJSON.Inclusions = inclusions
	accountJSON.Exclusions = exclusions

	err = b.updateAccount(ctx, req, name, accountJSON)
	if err != nil {
		return nil, err
	}
	_, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address":    account.Address.Hex(),
			"inclusions": accountJSON.Inclusions,
			"exclusions": accountJSON.Exclusions,
		},
	}, nil

}

func pathExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %v", err)
	}

	return out != nil, nil
}

// returns (nonce, toAddress, amount, gasPrice, gasLimit, error)

func (b *PluginBackend) getData(client *ethclient.Client, fromAddress common.Address, data *framework.FieldData) (*TransactionParams, error) {
	transactionParams, err := b.getBaseData(client, fromAddress, data, "to")
	if err != nil {
		return nil, err
	}
	var gasLimitIn *big.Int

	gasLimitField, ok := data.GetOk("gas_limit")

	if ok {
		gasLimitIn = util.ValidNumber(gasLimitField.(string))
	} else {
		// if no gaslimit is provided, take default or zero value
		gasLimitIn = util.ValidNumber(data.GetDefaultOrZero("gas_limit").(string))
	}

	gasLimit := gasLimitIn.Uint64()

	return &TransactionParams{
		Nonce:    transactionParams.Nonce,
		Address:  transactionParams.Address,
		Amount:   transactionParams.Amount,
		GasPrice: transactionParams.GasPrice,
		GasLimit: gasLimit,
	}, nil
}

// NewWalletTransactor is used with Token contracts
func (b *PluginBackend) NewWalletTransactor(chainID *big.Int, hdwallet *bip44.Wallet, account *accounts.Account) (*bind.TransactOpts, error) {
	return &bind.TransactOpts{
		From: account.Address,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != account.Address {
				return nil, errors.New("not authorized to sign this account")
			}
			signedTx, err := hdwallet.SignTx(*account, tx, chainID)
			if err != nil {
				return nil, err
			}

			return signedTx, nil
		},
	}, nil
}

func (b *PluginBackend) getBaseData(client *ethclient.Client, fromAddress common.Address, data *framework.FieldData, addressField string) (*TransactionParams, error) {
	var err error
	var address common.Address
	nonceData := "0"
	var nonce uint64
	var amount *big.Int
	var gasPriceIn *big.Int
	_, ok := data.GetOk("amount")
	if ok {
		amount = util.ValidNumber(data.Get("amount").(string))
		if amount == nil {
			return nil, fmt.Errorf("invalid amount")
		}
	} else {
		amount = util.ValidNumber("0")
	}

	_, ok = data.GetOk("nonce")
	if ok {
		nonceData = data.Get("nonce").(string)
		nonceIn := util.ValidNumber(nonceData)
		nonce = nonceIn.Uint64()
	} else {
		nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return nil, err
		}
	}

	_, ok = data.GetOk("gas_price")
	if ok {
		gasPriceIn = util.ValidNumber(data.Get("gas_price").(string))
		if gasPriceIn == nil {
			return nil, fmt.Errorf("invalid gas price")
		}
	} else {
		gasPriceIn = util.ValidNumber("0")
	}

	if big.NewInt(0).Cmp(gasPriceIn) == 0 {
		gasPriceIn, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, err
		}
	}

	if addressField != Empty {
		address = common.HexToAddress(data.Get(addressField).(string))
		return &TransactionParams{
			Nonce:    nonce,
			Address:  &address,
			Amount:   amount,
			GasPrice: gasPriceIn,
			GasLimit: 0,
		}, nil
	}
	return &TransactionParams{
		Nonce:    nonce,
		Address:  nil,
		Amount:   amount,
		GasPrice: gasPriceIn,
		GasLimit: 0,
	}, nil

}

func (b *PluginBackend) pathTransfer(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var txDataToSign []byte
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}
	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to " + config.getRPCURL())
	}

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}

	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	transactionParams, err := b.getData(client, account.Address, data)

	if err != nil {
		return nil, err
	}
	accountJSON.Inclusions = append(accountJSON.Inclusions, config.Inclusions...)
	accountJSON.Inclusions = append(accountJSON.Inclusions, accountJSON.Inclusions...)
	if len(accountJSON.Inclusions) > 0 && !util.Contains(accountJSON.Inclusions, transactionParams.Address.Hex()) {
		return nil, fmt.Errorf("%s violates the inclusions %+v", transactionParams.Address.Hex(), accountJSON.Inclusions)
	}
	err = config.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}
	err = accountJSON.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(transactionParams.Nonce, *transactionParams.Address, transactionParams.Amount, transactionParams.GasLimit, transactionParams.GasPrice, txDataToSign)
	signedTx, err := wallet.SignTx(*account, tx, chainID)
	if err != nil {
		return nil, err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	var signedTxBuff bytes.Buffer
	signedTx.EncodeRLP(&signedTxBuff)

	return &logical.Response{
		Data: map[string]interface{}{
			"transaction_hash":   signedTx.Hash().Hex(),
			"signed_transaction": hexutil.Encode(signedTxBuff.Bytes()),
			"from":               account.Address.Hex(),
			"to":                 transactionParams.Address.String(),
			"amount":             transactionParams.Amount.String(),
			"nonce":              strconv.FormatUint(transactionParams.Nonce, 10),
			"gas_price":          transactionParams.GasPrice.String(),
			"gas_limit":          strconv.FormatUint(transactionParams.GasLimit, 10),
		},
	}, nil
}

func (b *PluginBackend) pathDeploy(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}
	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to " + config.getRPCURL())
	}

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}

	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	transactionParams, err := b.getBaseData(client, account.Address, data, Empty)

	if err != nil {
		return nil, err
	}

	abiData := data.Get("abi").(string)
	parsed, err := abi.JSON(strings.NewReader(abiData))
	if err != nil {
		return nil, err
	}
	binData := data.Get("bin").(string)
	if err != nil {
		return nil, err
	}
	binRaw := common.FromHex(binData)
	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}
	gasLimitIn := util.ValidNumber(data.Get("gas_limit").(string))
	gasLimit := gasLimitIn.Uint64()

	transactOpts.GasPrice = transactionParams.GasPrice
	transactOpts.Nonce = big.NewInt(int64(transactionParams.Nonce))
	transactOpts.Value = big.NewInt(0) // in wei

	gasLimit, err = util.EstimateGas(transactOpts, parsed, binRaw, client)
	if err != nil {
		return nil, err
	}
	transactOpts.GasLimit = gasLimit
	contractAddress, tx, _, err := bind.DeployContract(transactOpts, parsed, binRaw, client)
	if err != nil {
		return nil, err
	}
	//	b.LogTx(tx)
	var signedTxBuff bytes.Buffer
	tx.EncodeRLP(&signedTxBuff)

	return &logical.Response{
		Data: map[string]interface{}{
			"transaction_hash":   tx.Hash().Hex(),
			"signed_transaction": hexutil.Encode(signedTxBuff.Bytes()),
			"from":               account.Address.Hex(),
			"contract":           contractAddress.Hex(),
			"nonce":              transactOpts.Nonce.String(),
			"gas_price":          transactOpts.GasPrice.String(),
			"gas_limit":          strconv.FormatUint(gasLimit, 10),
		},
	}, nil
}

func (b *PluginBackend) pathSignTx(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var txDataToSign []byte
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to " + config.getRPCURL())
	}

	name := data.Get("name").(string)

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}
	dataOrFile := data.Get("data").(string)
	encoding := data.Get("encoding").(string)
	if encoding == "hex" {
		txDataToSign, err = util.Decode([]byte(dataOrFile))
		if err != nil {
			return nil, err
		}
	} else if encoding == "utf8" {
		txDataToSign = []byte(dataOrFile)
	} else {
		return nil, fmt.Errorf("invalid encoding encountered - %s", encoding)
	}
	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}

	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}
	transactionParams, err := b.getData(client, account.Address, data)
	if err != nil {
		return nil, err
	}

	accountJSON.Inclusions = append(accountJSON.Inclusions, config.Inclusions...)
	if len(accountJSON.Inclusions) > 0 && !util.Contains(accountJSON.Inclusions, transactionParams.Address.Hex()) {
		return nil, fmt.Errorf("%s violates the set of inclusions %+v", transactionParams.Address.Hex(), accountJSON.Inclusions)
	}
	err = config.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}
	err = accountJSON.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(transactionParams.Nonce, *transactionParams.Address, transactionParams.Amount, transactionParams.GasLimit, transactionParams.GasPrice, txDataToSign)

	signedTx, err := wallet.SignTx(*account, tx, chainID)
	if err != nil {
		return nil, err
	}
	var signedTxBuff bytes.Buffer
	signedTx.EncodeRLP(&signedTxBuff)

	return &logical.Response{
		Data: map[string]interface{}{
			"transaction_hash":   signedTx.Hash().Hex(),
			"signed_transaction": hexutil.Encode(signedTxBuff.Bytes()),
			"from":               account.Address.Hex(),
			"to":                 transactionParams.Address.String(),
			"amount":             transactionParams.Amount.String(),
			"nonce":              strconv.FormatUint(transactionParams.Nonce, 10),
			"gas_price":          transactionParams.GasPrice.String(),
			"gas_limit":          strconv.FormatUint(transactionParams.GasLimit, 10),
		},
	}, nil

}

func (b *PluginBackend) pathReadBalance(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)
	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}

	_, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}
	balance, err := client.BalanceAt(context.Background(), account.Address, nil)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address": account.Address.Hex(),
			"balance": balance.String(),
		},
	}, nil

}

// LogTx is for debugging
func (b *PluginBackend) LogTx(tx *types.Transaction) {
	b.Logger().Info(fmt.Sprintf("\nTX DATA: %s\nGAS: %d\nGAS PRICE: %d\nVALUE: %d\nNONCE: %d\nTO: %s\n", hexutil.Encode(tx.Data()), tx.Gas(), tx.GasPrice(), tx.Value(), tx.Nonce(), tx.To().Hex()))
}

func (b *PluginBackend) pathSignMessage(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	message := data.Get("message").(string)
	name := data.Get("name").(string)

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}

	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	hashedMessage, _ := accounts.TextAndHash([]byte(message))

	signedMessage, err := wallet.SignHash(*account, []byte(hashedMessage))
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature":     hexutil.Encode(signedMessage),
			"address":       account.Address,
			"hashedMessage": hexutil.Encode(hashedMessage),
		},
	}, nil
}
