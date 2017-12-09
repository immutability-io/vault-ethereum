package ethereum

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func exportPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:      "export/" + framework.GenericNameRegex("name"),
			HelpSynopsis: "Export a single Ethereum JSON keystore from vault into. ",
			HelpDescription: `

Writes a JSON keystore to a folder (e.g., ~/.ethereum/keystore).

`,
			Fields: map[string]*framework.FieldSchema{
				"path": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Path to write to.",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathExportCreate,
			},
		},
	}
}

func (b *backend) pathExportCreate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("pathExportCreate", "path", req.Path)
	return nil, nil
}
