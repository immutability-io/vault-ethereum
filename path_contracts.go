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
	"errors"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/shopspring/decimal"
)

// Contract is the address of a contract
type Contract struct {
	Address         string `json:"address"`
	TransactionHash string `json:"transaction_hash"`
}

func contractsPaths(b *EthereumBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:      "deploy/" + framework.GenericNameRegex("name") + "/contracts/" + framework.GenericNameRegex("contract"),
			HelpSynopsis: "Sign and deploy an Ethereum contract.",
			HelpDescription: `

Deploys an Ethereum contract.

`,
			Fields: map[string]*framework.FieldSchema{
				"name":     &framework.FieldSchema{Type: framework.TypeString},
				"contract": &framework.FieldSchema{Type: framework.TypeString},
				"transaction_data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The transaction data.",
				},
				"amount": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Amount of ETH to fund the contract in Wei.",
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
					Description: "The gas limit in Wei for the transaction.",
				},
				"send": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "Send the transaction to the network.",
					Default:     true,
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathCreateContract,
				logical.ReadOperation:   b.pathReadContract,
				logical.DeleteOperation: b.pathContractsDelete,
			},
		},
	}
}

// NewTransactor is for contract deployment
func (b *EthereumBackend) NewTransactor(key *ecdsa.PrivateKey) *bind.TransactOpts {
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

func (b *EthereumBackend) pathCreateContract(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	input := []byte(data.Get("transaction_data").(string))
	name := data.Get("name").(string)
	sendTransaction := data.Get("send").(bool)
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("error reading account")
	}
	if account == nil {
		return nil, nil
	}
	toAddress := common.HexToAddress(account.Address)
	balance, _, _, err := b.readAccountBalance(ctx, req, name)
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
	if valid, err := b.validAccountConstraints(account, amount, account.Address); !valid {
		return nil, err
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
	gasLimitIn := ValidNumber(data.Get("gas_limit").(string))
	if gasLimitIn == nil {
		return nil, fmt.Errorf("invalid gas limit")
	}
	gasLimit := gasLimitIn.Uint64()
	if big.NewInt(0).Cmp(gasLimitIn) == 0 {
		gasLimit, err = client.EstimateGas(context.Background(), ethereum.CallMsg{
			To:   &toAddress,
			Data: input,
		})
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

	transactor := b.NewTransactor(privateKey)
	var rawTx *types.Transaction
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	rawTx = types.NewContractCreation(nonce, amount, gasLimit, gasPrice, input)

	signedTx, err := transactor.Signer(types.NewEIP155Signer(chainID), toAddress, rawTx)
	if err != nil {
		return nil, err
	}
	if sendTransaction {
		err = client.SendTransaction(ctx, signedTx)
		if err != nil {
			return nil, err
		}
	}

	contractJSON := &Contract{TransactionHash: signedTx.Hash().Hex()}
	entry, err := logical.StorageEntryJSON(req.Path, contractJSON)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}
	totalSpend, err := b.updateTotalSpend(ctx, req, fmt.Sprintf("accounts/%s", name), account, amount)
	if err != nil {
		return nil, err
	}
	amountInUSD, _ := decimal.NewFromString("0")
	if config.ChainID == EthereumMainnet {
		amountInUSD, err = ConvertToUSD(amount.String(), config.CoinMarketCapAPIKey)
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
			"total_spend":        totalSpend,
			"amount":             amount.String(),
			"amount_in_usd":      amountInUSD,
		},
	}, nil

}

func (b *EthereumBackend) pathReadContract(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	entry, err := req.Storage.Get(ctx, req.Path)
	var contract Contract
	err = entry.DecodeJSON(&contract)

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	name := data.Get("name").(string)
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("error reading account")
	}
	if account == nil {
		return nil, nil
	}

	hash := common.HexToHash(contract.TransactionHash)

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to " + config.getRPCURL())
	}

	receipt, err := client.TransactionReceipt(context.Background(), hash)
	var receiptAddress string
	if err != nil {
		receiptAddress = "receipt not available"
	} else {
		receiptAddress = receipt.ContractAddress.Hex()
	}

	return &logical.Response{
		Data: map[string]interface{}{

			"transaction_hash": contract.TransactionHash,
			"address":          receiptAddress,
		},
	}, nil
}

func (b *EthereumBackend) pathContractsDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Delete(ctx, req.Path); err != nil {
		return nil, err
	}

	return nil, nil
}
