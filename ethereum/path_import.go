package ethereum

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	MAX_KEYSTORE_SIZE int64  = 1024 // Just a heuristic to prevent reading stupid big files
	PATH_IMPORT       string = "import"
	PATH_ACCOUNTS     string = "accounts"
)

// accountsPaths is used to test CRUD and List operations. It is a simplified
// version of the passthrough backend that only accepts string values.
func importPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:      "import/" + framework.GenericNameRegex("name"),
			HelpSynopsis: "Generate and return a random password",
			HelpDescription: `

Generates a high-entropy password with the provided length and requirements,
returning it as part of the response. The generated password is not stored.

`,
			Fields: map[string]*framework.FieldSchema{
				"path":       &framework.FieldSchema{Type: framework.TypeString},
				"passphrase": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathImportCreate,
			},
		},
	}
}

func zeroKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

func (b *backend) importJSONKeystore(keystorePath string, passphrase string) (string, []byte, error) {
	var key *keystore.Key
	var jsonKeystore []byte
	file, err := os.Open(keystorePath)
	defer file.Close()
	stat, _ := file.Stat()
	if stat.Size() > MAX_KEYSTORE_SIZE {
		err = fmt.Errorf("keystore is suspiciously large at %d bytes", stat.Size())
	} else {
		jsonKeystore, err = ioutil.ReadFile(keystorePath)
		if err != nil {
			return "", nil, err
		}
		key, err = keystore.DecryptKey(jsonKeystore, passphrase)

		if key != nil && key.PrivateKey != nil {
			defer zeroKey(key.PrivateKey)
		}
	}
	return key.Address.Hex(), jsonKeystore, err
}

func pathExists(req *logical.Request, path string) (bool, error) {
	out, err := req.Storage.Get(path)
	if err != nil {
		return false, fmt.Errorf("existence check failed for %s: %v", path, err)
	}

	return out != nil, nil
}

func (b *backend) pathImportExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	accountPath := strings.Replace(req.Path, PATH_IMPORT, PATH_ACCOUNTS, -1)
	return pathExists(req, accountPath)
}

func (b *backend) pathImportCreate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	accountPath := strings.Replace(req.Path, PATH_IMPORT, PATH_ACCOUNTS, -1)
	exists, err := pathExists(req, accountPath)
	if !exists || err != nil {
		keystorePath := data.Get("path").(string)
		passphrase := data.Get("passphrase").(string)
		address, jsonKeystore, err := b.importJSONKeystore(keystorePath, passphrase)
		if err != nil {
			return nil, err
		}

		account := &Account{Address: address, Passphrase: passphrase, JSONKeystore: jsonKeystore}
		entry, err := logical.StorageEntryJSON(accountPath, account)

		err = req.Storage.Put(entry)
		if err != nil {
			return nil, err
		}
		return &logical.Response{
			Data: map[string]interface{}{
				"address":       address,
				"json_keystore": jsonKeystore,
			},
		}, nil
	} else {
		return nil, fmt.Errorf("this path %s exists. You can't import on top of it", accountPath)
	}
}
