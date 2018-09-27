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

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type AccountNames struct {
	Names []string `json:"names"`
}

func namesPaths(b *EthereumBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "names/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathNamesList,
			},
			HelpSynopsis: "List all the account names",
			HelpDescription: `
			All the names of accounts will be listed.
			`,
		},
		&framework.Path{
			Pattern:      "names/" + framework.GenericNameRegex("name"),
			HelpSynopsis: "Lookup a account's address by name.",
			HelpDescription: `

			Lookup a account's address by name.
`,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathNamesRead,
			},
		},
		&framework.Path{
			Pattern:      "names/" + framework.GenericNameRegex("name") + "/balance",
			HelpSynopsis: "Retrieve this accounts balance.",
			HelpDescription: `

			Retrieve this accounts balance.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathAccountBalanceRead,
			},
		},
		&framework.Path{
			Pattern:      "names/" + framework.GenericNameRegex("name") + "/verify",
			HelpSynopsis: "Verify that data was signed by a particular named account.",
			HelpDescription: `

			Verify that data was signed by a particular named account
`,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
				"data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The data to verify the signature of.",
				},
				"signature": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The signature to verify.",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathNamesVerify,
			},
		},
	}
}

// swagger:route  GET /{mount-path}/names/{name} Names pathNamesRead
//
// Handler returning the list of addresses associated with a named account
//
// ### This endpoint will list the addresses associated with a named account.
//
// ## Inputs:
//
// | Name    | Type     | Required | Default | Description                |
// | ------- | -------- | -------- | ---------| -------------------------- |
// | mount-path   | string    | true  | The endpoint configured for the plugin mount. |
// | name   | string    | true  | Specifies the name of the account to read. This is specified as part of the URL. |
//
// Responses:
//        200: AddressListResponse
func (b *EthereumBackend) pathNamesRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)
	account, err := b.readName(ctx, req, name)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, nil
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"address": account.Address,
		},
	}, nil
}

// swagger:route  GET /{mount-path}/names Names pathNamesList
//
// Handler returning the list of account names.
//
// ### This endpoint will list all account names.
//
// ## Inputs:
//
// | Name    | Type     | Required | Default | Description                |
// | ------- | -------- | -------- | ---------| -------------------------- |
// | mount-path   | string    | true  | The endpoint configured for the plugin mount. |
//
// Responses:
//        200: KeyListResponse
func (b *EthereumBackend) pathNamesList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	vals, err := req.Storage.List(ctx, "names/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

func (b *EthereumBackend) readName(ctx context.Context, req *logical.Request, name string) (*AccountAddress, error) {
	path := fmt.Sprintf("names/%s", name)
	entry, err := req.Storage.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var accountAddress AccountAddress
	err = entry.DecodeJSON(&accountAddress)

	if entry == nil {
		return nil, fmt.Errorf("failed to deserialize named account at %s", path)
	}

	return &accountAddress, nil
}

// swagger:route  POST /{mount-path}/names/{name}/verify Names pathNamesVerify
//
// Handler returning the list of account names.
//
// ### This endpoint will list all account names.
//
// ## Inputs:
//
// | Name    | Type     | Required | Default | Description                |
// | ------- | -------- | -------- | ---------| -------------------------- |
// | mount-path   | string    | true  | The endpoint configured for the plugin mount. |
// | name   | string    | true  | Specifies the name of the account to read. This is specified as part of the URL. |
// | data   | string    | true  | Some data. |
// | signature   | string    | true  | The signature to verify. |
//
// Responses:
//        200: VerifiedResponse
func (b *EthereumBackend) pathNamesVerify(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	name := data.Get("name").(string)
	return b.verifySignature(ctx, req, data, name)
}

func (b *EthereumBackend) pathAccountBalanceRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	balance, address, exchangeValue, err := b.readAccountBalance(ctx, req, name)
	if err != nil {
		return nil, err
	}
	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"address":        address,
			"balance":        balance.String(),
			"balance_in_usd": exchangeValue,
		},
	}, nil

}
