package ethereum

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func importPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:      "import/" + framework.GenericNameRegex("name"),
			HelpSynopsis: "Import a single Ethereum JSON keystore. ",
			HelpDescription: `

Reads a JSON keystore, decrypts it and stores the passphrase.

`,
			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Path to the keystore file - not the parent directory.",
				},
				"passphrase": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Passphrase used to encrypt private key - will not be returned.",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathImportCreate,
			},
		},
	}
}

func (b *backend) pathImportExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	b.Logger().Info("pathImportExistenceCheck", "path", req.Path)
	accountPath := strings.Replace(req.Path, PATH_IMPORT, PATH_ACCOUNTS, -1)
	return pathExists(req, accountPath)
}

func (b *backend) pathImportCreate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathImportCreate", "path", req.Path)
	accountPath := strings.Replace(req.Path, PATH_IMPORT, PATH_ACCOUNTS, -1)
	exists, err := pathExists(req, accountPath)
	if !exists || err != nil {
		keystorePath := data.Get("path").(string)
		passphrase := data.Get("passphrase").(string)
		address, jsonKeystore, err := b.importJSONKeystore(keystorePath, passphrase)
		if err != nil {
			return nil, err
		}

		accountJSON := &Account{Address: address, Passphrase: passphrase, JSONKeystore: jsonKeystore}
		entry, err := logical.StorageEntryJSON(accountPath, accountJSON)

		err = req.Storage.Put(entry)
		if err != nil {
			return nil, err
		}
		return &logical.Response{
			Data: map[string]interface{}{
				"address": address,
			},
		}, nil
	} else {
		return nil, fmt.Errorf("this path %s exists. You can't import on top of it", accountPath)
	}
}
