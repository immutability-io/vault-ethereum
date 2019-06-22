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
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/sha3"
)

const (

	// MaxKeystoreSize is just a guess to prevent injection attacks
	MaxKeystoreSize int64 = 1024 // Just a heuristic to prevent reading stupid big files
	// PassphraseWords is for big passphrases
	PassphraseWords int = 9
	// PassphraseSeparator is how we separate words
	PassphraseSeparator string = "-"
)

func importPaths(b *EthereumBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:      "import/" + framework.GenericNameRegex("name"),
			HelpSynopsis: "Import a single Ethereum JSON keystore. ",
			HelpDescription: `

Reads a JSON keystore, decrypts it and stores the passphrase.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
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

func (b *EthereumBackend) readJSONKeystore(keystorePath string) ([]byte, error) {
	var jsonKeystore []byte
	file, err := os.Open(keystorePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Size() > MaxKeystoreSize {
		err = fmt.Errorf("keystore is suspiciously large at %d bytes", stat.Size())
		return nil, err
	}
	jsonKeystore, err = ioutil.ReadFile(keystorePath)
	if err != nil {
		return nil, err
	}
	return jsonKeystore, nil
}

func (b *EthereumBackend) importJSONKeystore(ctx context.Context, keystorePath string, passphrase string) (*ecdsa.PrivateKey, error) {
	var key *keystore.Key
	jsonKeystore, err := b.readJSONKeystore(keystorePath)
	if err != nil {
		return nil, err
	}
	key, err = keystore.DecryptKey(jsonKeystore, passphrase)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, fmt.Errorf("failed to decrypt key")
	}

	return key.PrivateKey, err
}

func (b *EthereumBackend) pathImportCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	account, err := b.readAccount(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("Error reading account")
	}
	if account == nil {
		keystorePath := data.Get("path").(string)
		passphrase := data.Get("passphrase").(string)
		privateKey, err := b.importJSONKeystore(ctx, keystorePath, passphrase)
		if err != nil {
			return nil, err
		}
		defer ZeroKey(privateKey)
		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyString := hexutil.Encode(privateKeyBytes)[2:]

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("error casting public key to ECDSA")
		}

		publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
		publicKeyString := hexutil.Encode(publicKeyBytes)[4:]

		hash := sha3.NewLegacyKeccak256()
		hash.Write(publicKeyBytes[1:])
		address := hexutil.Encode(hash.Sum(nil)[12:])

		accountJSON := &Account{
			Address:    address,
			PrivateKey: privateKeyString,
			PublicKey:  publicKeyString,
		}
		path := fmt.Sprintf("accounts/%s", name)
		entry, err := logical.StorageEntryJSON(path, accountJSON)
		if err != nil {
			return nil, err
		}

		err = req.Storage.Put(ctx, entry)
		if err != nil {
			return nil, err
		}
		b.crossReference(ctx, req, name, accountJSON.Address)
		return &logical.Response{
			Data: map[string]interface{}{
				"address":              accountJSON.Address,
				"whitelist":            accountJSON.Whitelist,
				"blacklist":            accountJSON.Blacklist,
				"spending_limit_tx":    accountJSON.SpendingLimitTx,
				"spending_limit_total": accountJSON.SpendingLimitTotal,
				"total_spend":          accountJSON.TotalSpend,
			},
		}, nil
	}

	return nil, fmt.Errorf("account %s exists", name)
}
