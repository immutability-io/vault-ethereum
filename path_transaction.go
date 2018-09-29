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
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func transactionPaths(b *EthereumBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:      "transaction/" + framework.GenericNameRegex("id"),
			HelpSynopsis: "Query a transaction on the blockchain. ",
			HelpDescription: `

Query a transaction on the blockchain.

`,
			Fields: map[string]*framework.FieldSchema{
				"id": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathTransactionRead,
			},
		},
	}
}

func (b *EthereumBackend) pathTransactionRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	chainID := ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	id := data.Get("id").(string)
	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to " + config.getRPCURL())
	}
	txHash := common.HexToHash(id)
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, nil
	}

	msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
	if err != nil {
		return nil, err
	}

	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"transaction_hash": tx.Hash().Hex(),
			"pending":          isPending,
			"value":            tx.Value().String(),
			"gas":              tx.Gas(),
			"gas_price":        tx.GasPrice().Uint64(),
			"nonce":            tx.Nonce(),
			"address_to":       tx.To().Hex(),
			"address_from":     msg.From().Hex(),
			"receipt_status":   receipt.Status,
		},
	}, nil
}
