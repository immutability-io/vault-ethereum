package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// AccountAddress stores the name of the account to allow reverse lookup by address
type AccountAddress struct {
	Address string `json:"address"`
}

func addressesPaths(b *EthereumBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "addresses/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathAddressesList,
			},
			HelpSynopsis: "List all the account addresses",
			HelpDescription: `
			All the addresses of accounts will be listed.
			`,
		},
		&framework.Path{
			Pattern:      "addresses/" + framework.GenericNameRegex("address"),
			HelpSynopsis: "Lookup a account's name by address.",
			HelpDescription: `

			Lookup a account's name by address.
`,
			Fields: map[string]*framework.FieldSchema{
				"address": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathAddressesRead,
			},
		},
		&framework.Path{
			Pattern:      "addresses/" + framework.GenericNameRegex("address") + "/balance",
			HelpSynopsis: "Retrieve this accounts balance.",
			HelpDescription: `

			Retrieve this accounts balance.

`,
			Fields: map[string]*framework.FieldSchema{
				"address": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathAccountBalanceReadByAddress,
			},
		},
		&framework.Path{
			Pattern:      "addresses/" + framework.GenericNameRegex("address") + "/verify",
			HelpSynopsis: "Verify that data was signed by a particular address.",
			HelpDescription: `

			Verify that data was signed by a particular address
`,
			Fields: map[string]*framework.FieldSchema{
				"address": &framework.FieldSchema{Type: framework.TypeString},
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
				logical.UpdateOperation: b.pathAddressesVerify,
			},
		},
	}
}

// swagger:route GET /addresses/{address} Addresses pathAddressesRead
//
// Handler returning Account Names for an Address.
//
// Responses:
//        200: AccountNamesResponse
func (b *EthereumBackend) pathAddressesRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	address := data.Get("address").(string)
	account, err := b.readAddress(ctx, req, address)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, nil
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"names": account.Names,
		},
	}, nil
}

// swagger:route  GET /addresses Addresses pathAddressesList
//
// Handler returning the list of addresses.
//
// Responses:
//        200: AddressesResponse
func (b *EthereumBackend) pathAddressesList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	vals, err := req.Storage.List(ctx, "addresses/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

func (b *EthereumBackend) readAddress(ctx context.Context, req *logical.Request, address string) (*AccountNames, error) {
	path := fmt.Sprintf("addresses/%s", address)
	entry, err := req.Storage.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var accountNames AccountNames
	err = entry.DecodeJSON(&accountNames)
	if entry == nil {
		return nil, fmt.Errorf("failed to deserialize account at %s", path)
	}

	return &accountNames, nil
}

// swagger:route  POST /addresses/{address}/verify Addresses pathAddressesVerify
//
// Handler verifying that this account signed some data.
//
// Responses:
//        200: AddressesVerifiedResponse
func (b *EthereumBackend) pathAddressesVerify(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	address := data.Get("address").(string)
	account, err := b.readAddress(ctx, req, address)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, nil
	}
	if len(account.Names) == 0 {
		return nil, nil
	}

	return b.verifySignature(ctx, req, data, account.Names[0])
}

func (b *EthereumBackend) pathAccountBalanceReadByAddress(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	address := data.Get("address").(string)
	account, err := b.readAddress(ctx, req, address)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, nil
	}
	if len(account.Names) == 0 {
		return nil, nil
	}

	balance, _, err := b.readAccountBalance(ctx, req, account.Names[0])
	if err != nil {
		return nil, err
	}
	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"address": address,
			"balance": balance,
		},
	}, nil
}

func (b *EthereumBackend) crossReference(ctx context.Context, req *logical.Request, name, address string) error {
	accountAddress := &AccountAddress{Address: address}
	accountNames, err := b.readAddress(ctx, req, address)

	if accountNames == nil {
		accountNames = &AccountNames{}
	}
	accountNames.Names = append(accountNames.Names, name)

	pathAccountAddress := fmt.Sprintf("addresses/%s", accountAddress.Address)
	pathAccountName := fmt.Sprintf("names/%s", name)

	lookupNameEntry, err := logical.StorageEntryJSON(pathAccountName, accountAddress)
	if err != nil {
		return err
	}
	lookupAddressEntry, err := logical.StorageEntryJSON(pathAccountAddress, accountNames)

	if err != nil {
		return err
	}
	err = req.Storage.Put(ctx, lookupNameEntry)
	if err != nil {
		return err
	}
	err = req.Storage.Put(ctx, lookupAddressEntry)
	if err != nil {
		return err
	}

	return nil
}

func (b *EthereumBackend) removeCrossReference(ctx context.Context, req *logical.Request, name, address string) error {
	pathAccountAddress := fmt.Sprintf("addresses/%s", address)
	pathAccountName := fmt.Sprintf("names/%s", name)

	accountNames, err := b.readAddress(ctx, req, address)
	if err != nil {
		return err
	}
	if accountNames == nil || len(accountNames.Names) <= 1 {
		if err := req.Storage.Delete(ctx, pathAccountAddress); err != nil {
			return err
		}
	} else {
		updatedAccountNames := &AccountNames{}
		for i, accountName := range accountNames.Names {
			if accountName != name {
				updatedAccountNames.Names = append(updatedAccountNames.Names, accountNames.Names[i])
			}
		}
		lookupAddressEntry, err := logical.StorageEntryJSON(pathAccountAddress, updatedAccountNames)

		if err != nil {
			return err
		}
		err = req.Storage.Put(ctx, lookupAddressEntry)
		if err != nil {
			return err
		}
	}

	if err := req.Storage.Delete(ctx, pathAccountName); err != nil {
		return err
	}
	return nil
}
