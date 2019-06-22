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

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Transaction is an Ethereum transaction
type Transaction struct {
	Value     string `json:"value"`
	Gas       uint64 `json:"gas"`
	GasPrice  uint64 `json:"gas_price"`
	Nonce     uint64 `json:"nonce"`
	AddressTo string `json:"address_to"`
}

func blockPaths(b *EthereumBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:      "block/" + framework.GenericNameRegex("number"),
			HelpSynopsis: "Query a block on the blockchain. ",
			HelpDescription: `

Query a block on the blockchain.

`,
			Fields: map[string]*framework.FieldSchema{
				"number": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathBlockRead,
			},
		},
		&framework.Path{
			Pattern:      "block/" + framework.GenericNameRegex("number") + "/transactions",
			HelpSynopsis: "Get a list of all the transaction IDs on a block",
			HelpDescription: `

Get a list of all the transaction IDs on a block.

`,
			Fields: map[string]*framework.FieldSchema{
				"number": &framework.FieldSchema{Type: framework.TypeString},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathBlockTransactionsList,
			},
		},
	}
}

func (b *EthereumBackend) pathBlockRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	number := data.Get("number").(string)
	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to " + config.getRPCURL())
	}

	blockNumber := ValidNumber(number)

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"block":             block.Number().Uint64(),
			"time":              block.Time(),
			"difficulty":        block.Difficulty().Uint64(),
			"block_hash":        block.Hash().Hex(),
			"transaction_count": len(block.Transactions()),
		},
	}, nil
}

func (b *EthereumBackend) pathBlockTransactionsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	number := data.Get("number").(string)
	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to " + config.getRPCURL())
	}

	blockNumber := ValidNumber(number)

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, nil
	}

	response := make(map[string]interface{})
	for _, tx := range block.Transactions() {
		response[tx.Hash().Hex()] = &Transaction{
			Value:     tx.Value().String(),
			Gas:       tx.Gas(),
			GasPrice:  tx.GasPrice().Uint64(),
			Nonce:     tx.Nonce(),
			AddressTo: tx.To().Hex(),
		}
	}
	return &logical.Response{
		Data: response,
	}, nil

}
