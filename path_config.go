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
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/cidrutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	// EthereumMainnet Chain ID
	EthereumMainnet string = "1"
	// Morden Chain ID
	Morden string = "2"
	// Ropsten Chain ID
	Ropsten string = "3"
	// Rinkeby Chain ID
	Rinkeby string = "4"
	// RootstockMainnet Chain ID
	RootstockMainnet string = "30"
	// RootstockTestnet Chain ID
	RootstockTestnet string = "31"
	// Kovan Chain ID
	Kovan string = "42"
	// EthereumClassicMainnet Chain ID
	EthereumClassicMainnet string = "61"
	// EthereumClassicTestnet Chain ID
	EthereumClassicTestnet string = "62"
	// GethPrivateChains Chain ID
	GethPrivateChains string = "1337"
)

const (
	// InfuraMainnet is the default for EthereumMainnet
	InfuraMainnet string = "https://mainnet.infura.io"
	// InfuraRopsten is the default for Ropsten
	InfuraRopsten string = "https://ropsten.infura.io"
	// InfuraKovan is the default for Kovan
	InfuraKovan string = "https://kovan.infura.io"
	// InfuraRinkeby is the default for Rinkeby
	InfuraRinkeby string = "https://rinkeby.infura.io"
	// Local is the default for localhost
	Local string = "http://localhost:8545"
)

// Config contains the configuration for each mount
type Config struct {
	BoundCIDRList []string `json:"bound_cidr_list_list" structs:"bound_cidr_list" mapstructure:"bound_cidr_list"`
	RPC           string   `json:"rpc_url"`
	InfuraAPIKey  string   `json:"api_key"`
	ChainID       string   `json:"chain_id"`
}

func configPaths(b *EthereumBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "config",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathCreateConfig,
				logical.UpdateOperation: b.pathCreateConfig,
				logical.ReadOperation:   b.pathReadConfig,
			},
			HelpSynopsis: "Configure the trustee plugin.",
			HelpDescription: `
			Configure the trustee plugin.
			`,
			Fields: map[string]*framework.FieldSchema{
				"chain_id": &framework.FieldSchema{
					Type: framework.TypeString,
					Description: `Ethereum network - can be one of the following values:

					1 - Ethereum mainnet
					2 - Morden (disused), Expanse mainnet
					3 - Ropsten
					4 - Rinkeby
					30 - Rootstock mainnet
					31 - Rootstock testnet
					42 - Kovan
					61 - Ethereum Classic mainnet
					62 - Ethereum Classic testnet
					1337 - Geth private chains (default)`,
					Default: Rinkeby,
				},
				"rpc_url": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: `The RPC address of the Ethereuem network.`,
				},
				"api_key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: `The Infura API Key.`,
				},
				"bound_cidr_list": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Comma separated string or list of CIDR blocks. If set, specifies the blocks of
IP addresses which can perform the login operation.`,
				},
			},
		},
	}
}

func (config *Config) getRPCURL() string {
	if config.InfuraAPIKey != "" {
		url := strings.TrimRight(config.RPC, "/")
		if isInfuraNetwork(url) {
			return fmt.Sprintf("%s/%s", url, config.InfuraAPIKey)
		}
	}
	return config.RPC
}

func isInfuraNetwork(url string) bool {
	switch url {
	case InfuraMainnet:
		return true
	case InfuraRopsten:
		return true
	case InfuraRinkeby:
		return true
	case InfuraKovan:
		return true
	}
	return false
}

func getDefaultNetwork(chainID string) string {
	switch chainID {
	case EthereumMainnet:
		return InfuraMainnet
	case Ropsten:
		return InfuraRopsten
	case Rinkeby:
		return InfuraRinkeby
	case Kovan:
		return InfuraKovan
	}
	return Local
}

func (b *EthereumBackend) pathCreateConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	rpcURL := data.Get("rpc_url").(string)
	apiKey := data.Get("api_key").(string)
	chainID := data.Get("chain_id").(string)
	if rpcURL == "" {
		rpcURL = getDefaultNetwork(chainID)
	}
	var boundCIDRList []string
	if boundCIDRListRaw, ok := data.GetOk("bound_cidr_list"); ok {
		boundCIDRList = boundCIDRListRaw.([]string)
	}
	configBundle := Config{
		BoundCIDRList: boundCIDRList,
		RPC:           rpcURL,
		ChainID:       chainID,
		InfuraAPIKey:  apiKey,
	}
	entry, err := logical.StorageEntryJSON("config", configBundle)

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"bound_cidr_list": configBundle.BoundCIDRList,
			"chain_id":        configBundle.ChainID,
			"api_key":         configBundle.InfuraAPIKey,
			"rpc_url":         configBundle.RPC,
		},
	}, nil
}

func (b *EthereumBackend) pathReadConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	configBundle, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if configBundle == nil {
		return nil, nil
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"bound_cidr_list": configBundle.BoundCIDRList,
			"chain_id":        configBundle.ChainID,
			"api_key":         configBundle.InfuraAPIKey,
			"rpc_url":         configBundle.RPC,
		},
	}, nil
}

// Config returns the configuration for this EthereumBackend.
func (b *EthereumBackend) readConfig(ctx context.Context, s logical.Storage) (*Config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf("the ethereum backend is not configured properly")
	}

	var result Config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, fmt.Errorf("error reading configuration: %s", err)
		}
	}

	return &result, nil
}

func (b *EthereumBackend) configured(ctx context.Context, req *logical.Request) (*Config, error) {
	config, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("backend not properly configured")
	}
	if validConnection, err := b.validIPConstraints(config, req); !validConnection {
		return nil, err
	}

	return config, nil
}

func (b *EthereumBackend) validIPConstraints(config *Config, req *logical.Request) (bool, error) {
	if len(config.BoundCIDRList) != 0 {
		if req.Connection == nil || req.Connection.RemoteAddr == "" {
			return false, fmt.Errorf("failed to get connection information")
		}

		belongs, err := cidrutil.IPBelongsToCIDRBlocksSlice(req.Connection.RemoteAddr, config.BoundCIDRList)
		if err != nil {
			return false, errwrap.Wrapf("failed to verify the CIDR restrictions set on the role: {{err}}", err)
		}
		if !belongs {
			return false, fmt.Errorf("source address %q unauthorized through CIDR restrictions on the role", req.Connection.RemoteAddr)
		}
	}
	return true, nil
}
