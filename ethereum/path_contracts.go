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
	Address string `json:"contract_address"`
	Hash    string `json:"tx_hash"`
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
				"amount": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Amount of ETH (in Gwei).",
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
					Description: "The gas limit (in Gwei) for the transaction.",
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

func (b *backend) pathContractsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vals, err := req.Storage.List(ctx, req.Path)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

func (b *backend) pathCreateContract(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	amount := math.MustParseBig256(data.Get("amount").(string))
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
	account, err := b.readAccount(ctx, req, accountPath, true)
	if err != nil {
		return nil, err
	}
	allowed, err := b.isDeployAllowed(account, amount)
	if !allowed {
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
	rawTx = types.NewContractCreation(nonce, amount, gasLimit, gasPrice, input)

	signedTx, err := transactor.Signer(types.NewEIP155Signer(chainID), common.HexToAddress(account.Address), rawTx)
	if err != nil {
		return nil, err
	}
	err = ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	contractJSON := &Contract{Hash: signedTx.Hash().Hex()}
	entry, err := logical.StorageEntryJSON(req.Path, contractJSON)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"tx_hash": signedTx.Hash().Hex(),
		},
	}, nil
}

func (b *backend) pathReadContract(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get(ctx, req.Path)
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
	account, err := b.readAccount(ctx, req, accountPath, false)
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
		receiptAddress = "Receipt not available"
	} else {
		receiptAddress = receipt.ContractAddress.Hex()
	}

	return &logical.Response{
		Data: map[string]interface{}{

			"tx_hash": contract.Hash,
			"address": receiptAddress,
		},
	}, nil
}
