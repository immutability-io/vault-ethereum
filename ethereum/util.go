package ethereum

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	accounts "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	rpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/hashicorp/vault/logical"
)

const (
	ProtocolKeystore    string = "keystore://"
	MaxKeystoreSize     int64  = 1024 // Just a heuristic to prevent reading stupid big files
	RequestPathImport   string = "import"
	RequestPathAccounts string = "accounts"
)

func (b *backend) writeTemporaryKeystoreFile(path string, filename string, data []byte) (string, error) {
	keystorePath := path + "/" + filename
	err := ioutil.WriteFile(keystorePath, data, 0644)
	return keystorePath, err
}

func (b *backend) createTemporaryKeystoreDirectory() (string, error) {
	dir, err := ioutil.TempDir("", "keystore")
	return dir, err
}

func (b *backend) removeTemporaryKeystore(path string) error {
	return os.RemoveAll(path)
}

func convertMapToStringValue(initial map[string]interface{}) map[string]string {
	result := map[string]string{}
	for key, value := range initial {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}

func parseURL(url string) (accounts.URL, error) {
	parts := strings.Split(url, "://")
	if len(parts) != 2 || parts[0] == "" {
		return accounts.URL{}, errors.New("protocol scheme missing")
	}
	return accounts.URL{
		Scheme: parts[0],
		Path:   parts[1],
	}, nil
}

func (b *backend) rekeyJSONKeystore(keystorePath string, passphrase string, newPassphrase string) ([]byte, error) {
	var key *keystore.Key
	jsonKeystore, err := b.readJSONKeystore(keystorePath)
	if err != nil {
		return nil, err
	}
	key, _ = keystore.DecryptKey(jsonKeystore, passphrase)

	if key != nil && key.PrivateKey != nil {
		defer zeroKey(key.PrivateKey)
	}
	jsonBytes, err := keystore.EncryptKey(key, newPassphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	return jsonBytes, err
}

func (b *backend) readKeyFromJSONKeystore(keystorePath string, passphrase string) (*keystore.Key, error) {
	var key *keystore.Key
	jsonKeystore, err := b.readJSONKeystore(keystorePath)
	if err != nil {
		return nil, err
	}
	key, _ = keystore.DecryptKey(jsonKeystore, passphrase)

	if key != nil && key.PrivateKey != nil {
		return key, nil
	} else {
		return nil, fmt.Errorf("failed to read key from keystore")
	}
}

func zeroKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

func (b *backend) importJSONKeystore(keystorePath string, passphrase string) (string, []byte, error) {
	b.Logger().Info("importJSONKeystore", "keystorePath", keystorePath)
	var key *keystore.Key
	jsonKeystore, err := b.readJSONKeystore(keystorePath)
	if err != nil {
		return "", nil, err
	}
	key, err = keystore.DecryptKey(jsonKeystore, passphrase)

	if key != nil && key.PrivateKey != nil {
		defer zeroKey(key.PrivateKey)
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

func (b *backend) readJSONKeystore(keystorePath string) ([]byte, error) {
	b.Logger().Info("readJSONKeystore", "keystorePath", keystorePath)
	var jsonKeystore []byte
	file, err := os.Open(keystorePath)
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Size() > MaxKeystoreSize {
		err = fmt.Errorf("keystore is suspiciously large at %d bytes", stat.Size())
		return nil, err
	} else {
		jsonKeystore, err = ioutil.ReadFile(keystorePath)
		if err != nil {
			return nil, err
		}
		return jsonKeystore, nil
	}
}

func (b *backend) NewTransactor(key *ecdsa.PrivateKey) *bind.TransactOpts {
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}

func (b *backend) getAccountPrivateKey(path string, account Account) (*keystore.Key, error) {
	tmpDir, err := b.createTemporaryKeystoreDirectory()
	if err != nil {
		return nil, err
	}

	keystorePath, err := b.writeTemporaryKeystoreFile(tmpDir, account.KeystoreName, account.JSONKeystore)
	if err != nil {
		return nil, err
	}
	key, err := b.readKeyFromJSONKeystore(keystorePath, account.Passphrase)
	if err != nil {
		return nil, err
	}
	err = b.removeTemporaryKeystore(tmpDir)
	return key, err
}

func (b *backend) exportKeystore(path string, account *Account) (string, error) {
	directory, err := b.writeTemporaryKeystoreFile(path, account.KeystoreName, account.JSONKeystore)
	return directory, err
}

func (b *backend) readAccount(req *logical.Request, path string, deep bool) (*Account, error) {
	entry, err := req.Storage.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to find account at %s", path)
	}
	var account Account
	err = entry.DecodeJSON(&account)

	if entry == nil {
		return nil, fmt.Errorf("failed to deserialize account at %s", path)
	}

	if deep {
		client, err := rpc.Dial(account.RPC)
		if err != nil {
			return nil, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ethClient := ethclient.NewClient(client)
		accountAddr := common.HexToAddress(account.Address)
		pendingBalance, err := ethClient.PendingBalanceAt(ctx, accountAddr)
		if err != nil {
			return nil, err
		}
		pendingNonce, err := ethClient.PendingNonceAt(ctx, accountAddr)
		if err != nil {
			return nil, err
		}
		pendingTxCount, err := ethClient.PendingTransactionCount(ctx)
		if err != nil {
			return nil, err
		}
		account.PendingNonce = pendingNonce
		account.PendingBalance = pendingBalance
		account.PendingTxCount = pendingTxCount
	}
	return &account, nil
}
