package ethereum

import (
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/sethvargo/go-diceware/diceware"
)

type Account struct {
	Address      string `json:"address"` // Ethereum account address derived from the key
	Passphrase   string `json:"passphrase"`
	JSONKeystore []byte `json:"json_keystore"`
}

// accountsPaths is used to test CRUD and List operations. It is a simplified
// version of the passthrough backend that only accepts string values.
func accountsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "accounts/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathAccountsList,
			},
		},
		&framework.Path{
			Pattern:      "accounts/" + framework.GenericNameRegex("account"),
			HelpSynopsis: "Create an Ethereum account using a generated or provided passphrase",
			HelpDescription: `

Generates a high-entropy password with the provided length and requirements,
returning it as part of the response. The generated password is not stored.

`,
			Fields: map[string]*framework.FieldSchema{
				"passphrase":          &framework.FieldSchema{Type: framework.TypeString},
				"generate_passphrase": &framework.FieldSchema{Type: framework.TypeBool},
				"words": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Number of words for the passphrase.",
					Default:     6,
				},
				"separator": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Character to separate words in passphrase.",
					Default:     "-",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.pathAccountsRead,
				logical.CreateOperation: b.pathAccountsCreate,
				logical.UpdateOperation: b.pathAccountsUpdate,
				logical.DeleteOperation: b.pathAccountsDelete,
			},
		},
	}
}

func (b *backend) pathAccountsRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get(req.Path)
	var account Account
	err = entry.DecodeJSON(&account)

	if err != nil {
		return nil, err
	}
	b.Logger().Info("What is stored here", "entry", entry)
	if entry == nil {
		return nil, nil
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"address":       account.Address,
			"json_keystore": account.JSONKeystore,
		},
	}, nil
}

func (b *backend) pathAccountsCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	passphrase := data.Get("passphrase").(string)
	generatePassphrase := data.Get("generate_passphrase").(bool)

	words := data.Get("words").(int)
	separator := data.Get("separator").(string)

	if generatePassphrase {
		list, _ := diceware.Generate(words)
		passphrase = strings.Join(list, separator)
	}
	entry := &logical.StorageEntry{
		Key:      req.Path,
		Value:    []byte(passphrase),
		SealWrap: true,
	}

	s := req.Storage
	err := s.Put(entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"passphrase": passphrase,
		},
	}, nil
}

func (b *backend) pathAccountsUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get(req.Path)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"address":    entry.Value,
			"passphrase": entry.Value,
		},
	}, nil
}

func (b *backend) pathAccountsDelete(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(req.Path); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathAccountsList(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vals, err := req.Storage.List("accounts/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}
