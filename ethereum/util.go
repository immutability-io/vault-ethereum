package ethereum

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	accounts "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/hashicorp/vault/logical"
)

const (
	FS_TEMPORARY      string = "/tmp/"
	PROTOCOL_KEYSTORE string = "keystore://"
	MAX_KEYSTORE_SIZE int64  = 1024 // Just a heuristic to prevent reading stupid big files
	PATH_IMPORT       string = "import"
	PATH_ACCOUNTS     string = "accounts"
)

func (b *backend) buildKeystoreURL(filename string) string {
	return PROTOCOL_KEYSTORE + FS_TEMPORARY + filename
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
	stat, _ := file.Stat()
	if stat.Size() > MAX_KEYSTORE_SIZE {
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
