package ethereum

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	rpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type Contract struct {
	Address string `json:"address"` // Ethereum account address derived from the key
	Hash    string `json:"hash"`    // Ethereum account address derived from the key
}

func contractsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "accounts/" + framework.GenericNameRegex("name") + "/contracts/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathContractsList,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("name") + "/contracts/" + framework.GenericNameRegex("contract"),
			HelpSynopsis: "Sign and deploy an Ethereum contract.",
			HelpDescription: `

Deploys an Ethereum contract.

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
					Default:     "20000000000",
				},
				"gas_limit": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The gas limit for the transaction.",
					Default:     "50000",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathCreateContract,
				logical.ReadOperation:   b.pathReadContract,
			},
		},
	}
}

func (b *backend) pathContractsList(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathContractsList", "path", req.Path)
	vals, err := req.Storage.List(req.Path)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

// func (b *backend) createContract(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
//
// 	value := math.MustParseBig256(data.Get("value").(string))
// 	nonce := math.MustParseUint64(data.Get("nonce").(string))
// 	gasPrice := math.MustParseBig256(data.Get("gas_price").(string))
// 	gasLimit := math.MustParseBig256(data.Get("gas_limit").(string))
// 	input := []byte(data.Get("transaction_data").(string))
// 	chainID := math.MustParseBig256(account.ChainID)
//
// 	prunedPath := strings.Replace(req.Path, "/sign-contract", "", -1)
// 	account, err := b.readAccount(req, prunedPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	key, err := b.getAccountPrivateKey(prunedPath, *account)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer zeroKey(key.PrivateKey)
//
// 	transactor := b.NewTransactor(key.PrivateKey)
// 	var rawTx *types.Transaction
// 	rawTx = types.NewContractCreation(nonce, value, gasLimit, gasPrice, input)
// 	//contract := &Contract{Address: rawTx.Hash().Hex()}
// 	signedTx, err := transactor.Signer(types.NewEIP155Signer(chainID), common.HexToAddress(account.Address), rawTx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	encoded, err := rlp.EncodeToBytes(signedTx)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &logical.Response{
// 		Data: map[string]interface{}{
// 			"signed_tx": hexutil.Encode(encoded[:]),
// 		},
// 	}, nil
// }
//
func (b *backend) pathCreateContract(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathCreateContract", "path", req.Path)

	value := math.MustParseBig256(data.Get("value").(string))
	nonce := math.MustParseUint64(data.Get("nonce").(string))
	gasPrice := math.MustParseBig256(data.Get("gas_price").(string))
	gasLimit := math.MustParseBig256(data.Get("gas_limit").(string))
	input := []byte(data.Get("transaction_data").(string))
	var accountPath string
	parsedPath := strings.Split(req.Path, "/contracts/")
	if len(parsedPath) >= 1 {
		accountPath = parsedPath[0]
	} else {
		return nil, fmt.Errorf("something sketchy with the path: %s", req.Path)
	}
	account, err := b.readAccount(req, accountPath)
	if err != nil {
		return nil, err
	}
	chainID := math.MustParseBig256(account.ChainID)
	key, err := b.getAccountPrivateKey(accountPath, *account)
	if err != nil {
		return nil, err
	}
	defer zeroKey(key.PrivateKey)

	transactor := b.NewTransactor(key.PrivateKey)
	var rawTx *types.Transaction
	b.Logger().Info("Transaction Nonce", "nonce", nonce)
	client, err := rpc.Dial(account.RPC)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ethClient := ethclient.NewClient(client)
	fromAddress := common.HexToAddress(account.Address)
	nonce, err = ethClient.NonceAt(ctx, fromAddress, nil)
	if err != nil {
		return nil, err
	}
	rawTx = types.NewContractCreation(nonce, value, gasLimit, gasPrice, input)
	//contract := &Contract{Address: rawTx.Hash().Hex()}
	signedTx, err := transactor.Signer(types.NewEIP155Signer(chainID), common.HexToAddress(account.Address), rawTx)
	if err != nil {
		return nil, err
	}
	b.Logger().Info("Signed Transaction", "signedTx", signedTx.String())
	err = ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		b.Logger().Info("Error sending TX", "error", err)
		return nil, err
	} else {
		b.Logger().Info("Sent TX", "signedTx", signedTx.Hash().Hex())
	}
	contractJSON := &Contract{Hash: signedTx.Hash().Hex()}
	entry, err := logical.StorageEntryJSON(req.Path, contractJSON)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{

			"transaction_hash": signedTx.Hash().Hex(),
		},
	}, nil
}

func (b *backend) pathReadContract(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathReadContract", "path", req.Path)
	entry, err := req.Storage.Get(req.Path)
	var contract Contract
	err = entry.DecodeJSON(&contract)

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var accountPath string
	parsedPath := strings.Split(req.Path, "/contracts/")
	if len(parsedPath) >= 1 {
		accountPath = parsedPath[0]
	} else {
		return nil, fmt.Errorf("something sketchy with the path: %s", req.Path)
	}
	account, err := b.readAccount(req, accountPath)
	if err != nil {
		return nil, err
	}
	client, err := rpc.Dial(account.RPC)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ethClient := ethclient.NewClient(client)

	hash := common.HexToHash(contract.Hash)
	receipt, err := ethClient.TransactionReceipt(ctx, hash)
	var receiptAddress string
	if err != nil {
		b.Logger().Info("Error reading receipt", "err", fmt.Sprintf("%s", err))
		receiptAddress = "Receipt not available"
	} else {
		receiptAddress = receipt.ContractAddress.Hex()
	}

	return &logical.Response{
		Data: map[string]interface{}{

			"transaction_hash": contract.Hash,
			"contract_address": receiptAddress,
		},
	}, nil
}
