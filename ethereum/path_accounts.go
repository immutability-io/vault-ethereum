package ethereum

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/sethvargo/go-diceware/diceware"
)

const (
	FS_TEMPORARY      string = "/tmp/"
	PROTOCOL_KEYSTORE string = "keystore://"
)

type Account struct {
	Address      string `json:"address"` // Ethereum account address derived from the key
	Passphrase   string `json:"passphrase"`
	URL          string `json:"url"`
	JSONKeystore []byte `json:"json_keystore"`
}

func (b *backend) writeTemporaryKeystoreFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}

func (b *backend) createTemporaryKeystore(name string) (string, error) {
	file, _ := os.Open(FS_TEMPORARY + name)
	if file != nil {
		file.Close()
		return "", fmt.Errorf("account already exists at %s", FS_TEMPORARY+name)
	}
	return FS_TEMPORARY + name, os.MkdirAll(FS_TEMPORARY+name, os.FileMode(0522))
}

func (b *backend) removeTemporaryKeystore(name string) error {
	file, _ := os.Open(FS_TEMPORARY + name)
	if file != nil {
		return os.RemoveAll(FS_TEMPORARY + name)
	} else {
		return fmt.Errorf("keystore doesn't exist at %s", FS_TEMPORARY+name)
	}

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
			Pattern:      "accounts/" + framework.GenericNameRegex("name"),
			HelpSynopsis: "Create an Ethereum account using a generated or provided passphrase",
			HelpDescription: `

Creates (or updates) an Ethereum externally owned account (EOAs): an account controlled by a private key.
Optionally generates a high-entropy passphrase with the provided length and requirements. The passphrase
is not returned, but it is stored at a separate path (accounts/<name>/passphrase) to allow fine
grained access controls over exposure of the passphrase. The update operation will create a new keystore using
the new passphrase.

`,
			Fields: map[string]*framework.FieldSchema{
				"passphrase": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The passphrase used to encrypt the private key.",
				},
				"generate_passphrase": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "Generate the passphrase.",
					Default:     false,
				},
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

func (b *backend) pathAccountsRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get(req.Path)
	var account Account
	err = entry.DecodeJSON(&account)

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"address": account.Address,
		},
	}, nil
}

func (b *backend) pathAccountsCreate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	passphrase := data.Get("passphrase").(string)
	generatePassphrase := data.Get("generate_passphrase").(bool)
	words := data.Get("words").(int)
	separator := data.Get("separator").(string)

	if generatePassphrase {
		list, _ := diceware.Generate(words)
		passphrase = strings.Join(list, separator)
	}
	tmpDir, err := b.createTemporaryKeystore(req.Path)
	if err != nil {
		return nil, err
	}
	ks := keystore.NewKeyStore(tmpDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(passphrase)
	if err != nil {
		return nil, err
	}
	keystorePath := strings.Replace(account.URL.String(), PROTOCOL_KEYSTORE, "", -1)
	jsonKeystore, err := b.readJSONKeystore(keystorePath)
	accountJSON := &Account{Address: account.Address.Hex(), Passphrase: passphrase, URL: account.URL.String(), JSONKeystore: jsonKeystore}
	entry, err := logical.StorageEntryJSON(req.Path, accountJSON)
	b.Logger().Info("Create account", "passphrase", passphrase)
	b.Logger().Info("Create account", "url", account.URL.String())
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}
	b.removeTemporaryKeystore(req.Path)
	return &logical.Response{
		Data: map[string]interface{}{
			"account":  account.Address.Hex(),
			"keystore": fmt.Sprintf("%s", jsonKeystore),
		},
	}, nil
}

func (b *backend) rekeyJSONKeystore(keystorePath string, passphrase string, newPassphrase string) ([]byte, error) {
	var key *keystore.Key
	jsonKeystore, err := b.readJSONKeystore(keystorePath)
	if err != nil {
		return nil, err
	}
	key, err = keystore.DecryptKey(jsonKeystore, passphrase)

	if key != nil && key.PrivateKey != nil {
		defer zeroKey(key.PrivateKey)
	}
	jsonBytes, err := keystore.EncryptKey(key, newPassphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	return jsonBytes, err
}

func (b *backend) pathAccountsUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	passphrase := data.Get("passphrase").(string)
	generatePassphrase := data.Get("generate_passphrase").(bool)
	words := data.Get("words").(int)
	separator := data.Get("separator").(string)

	if generatePassphrase {
		list, _ := diceware.Generate(words)
		passphrase = strings.Join(list, separator)
	} else if passphrase == "" {
		return nil, fmt.Errorf("nothing to update - no passphrase supplied")
	}
	entry, err := req.Storage.Get(req.Path)
	if err != nil {
		return nil, err
	}
	var account Account
	err = entry.DecodeJSON(&account)

	if err != nil {
		return nil, err
	}
	_, err = b.createTemporaryKeystore(req.Path)
	if err != nil {
		return nil, err
	}
	keystorePath := strings.Replace(account.URL, PROTOCOL_KEYSTORE, "", -1)
	b.writeTemporaryKeystoreFile(keystorePath, account.JSONKeystore)

	jsonKeystore, err := b.rekeyJSONKeystore(keystorePath, account.Passphrase, passphrase)
	if err != nil {
		return nil, err
	} else {
		b.writeTemporaryKeystoreFile(keystorePath, jsonKeystore)
		account.Passphrase = passphrase
		account.JSONKeystore = jsonKeystore
		entry, _ = logical.StorageEntryJSON(req.Path, account)

		err = req.Storage.Put(entry)
		if err != nil {
			return nil, err
		}
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"address":  account.Address,
			"keystore": fmt.Sprintf("%s", account.JSONKeystore),
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
