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

	"github.com/immutability-io/vault-ethereum/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/logical"
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
	Kovan     string = "42"
	BinanceId        = "97"
	// EthereumClassicMainnet Chain ID
	EthereumClassicMainnet string = "61"
	// EthereumClassicTestnet Chain ID
	EthereumClassicTestnet string = "62"
	// GethPrivateChains Chain ID
	GethPrivateChains string = "1337"
	// InfuraMainnet is the default for EthereumMainnet
	InfuraMainnet string = "https://mainnet.infura.io"
	// InfuraRopsten is the default for Ropsten
	InfuraRopsten string = "https://ropsten.infura.io"
	// InfuraKovan is the default for Kovan
	InfuraKovan string = "https://kovan.infura.io"
	// InfuraRinkeby is the default for Rinkeby
	InfuraRinkeby string = "https://rinkeby.infura.io"
	// Local is the default for localhost
	Local     string = "http://localhost:8545"
	BinaceRpc        = "https://data-seed-prebsc-1-s1.binance.org:8545/"
)

// ConfigJSON contains the configuration for each mount
type ConfigJSON struct {
	BoundCIDRList []string `json:"bound_cidr_list_list" structs:"bound_cidr_list" mapstructure:"bound_cidr_list"`
	Inclusions    []string `json:"inclusions"`
	Exclusions    []string `json:"exclusions"`
	RPC           string   `json:"rpc_url"`
	ChainID       string   `json:"chain_id"`
}

// ValidAddress returns an error if the address is not included or if it is excluded
func (config *ConfigJSON) ValidAddress(toAddress *common.Address) error {
	if util.Contains(config.Exclusions, toAddress.Hex()) {
		return fmt.Errorf("%s is excludeded by this mount", toAddress.Hex())
	}

	if len(config.Inclusions) > 0 && !util.Contains(config.Inclusions, toAddress.Hex()) {
		return fmt.Errorf("%s is not in the set of inclusions of this mount", toAddress.Hex())
	}
	return nil
}

func configPaths(b *PluginBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: QualifiedPath("config"),
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWriteConfig,
				logical.UpdateOperation: b.pathWriteConfig,
				logical.ReadOperation:   b.pathReadConfig,
			},
			HelpSynopsis: "Configure the Vault Ethereum plugin.",
			HelpDescription: `
			Configure the Vault Ethereum plugin.
			`,
			Fields: map[string]*framework.FieldSchema{
				"chain_id": {
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
					Default: BinanceId,
				},
				"rpc_url": {
					Type:        framework.TypeString,
					Default:     BinaceRpc,
					Description: "The RPC address of the Ethereum network",
				},
				"inclusions": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Only these accounts may be transaction with",
				},
				"exclusions": {
					Type:        framework.TypeCommaStringSlice,
					Description: "These accounts can never be transacted with",
				},
				"bound_cidr_list": {
					Type: framework.TypeCommaStringSlice,
					Description: `Comma separated string or list of CIDR blocks.
If set, specifies the blocks of IPs which can perform the login operation;
if unset, there are no IP restrictions.`,
				},
			},
		},
	}
}

func (config *ConfigJSON) getRPCURL() string {
	return config.RPC
}

func (b *PluginBackend) pathWriteConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	rpcURL := data.Get("rpc_url").(string)
	chainID := data.Get("chain_id").(string)
	var boundCIDRList []string
	if boundCIDRListRaw, ok := data.GetOk("bound_cidr_list"); ok {
		boundCIDRList = boundCIDRListRaw.([]string)
	}
	var inclusions []string
	if inclusionsRaw, ok := data.GetOk("inclusions"); ok {
		inclusions = inclusionsRaw.([]string)
	}
	var exclusions []string
	if exclusionsRaw, ok := data.GetOk("exclusions"); ok {
		exclusions = exclusionsRaw.([]string)
	}
	configBundle := ConfigJSON{
		BoundCIDRList: boundCIDRList,
		Inclusions:    inclusions,
		Exclusions:    exclusions,
		ChainID:       chainID,
		RPC:           rpcURL,
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
			"inclusions":      configBundle.Inclusions,
			"exclusions":      configBundle.Exclusions,
			"rpc_url":         configBundle.RPC,
			"chain_id":        configBundle.ChainID,
		},
	}, nil
}

func (b *PluginBackend) pathReadConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
			"inclusions":      configBundle.Inclusions,
			"exclusions":      configBundle.Exclusions,
			"rpc_url":         configBundle.RPC,
			"chain_id":        configBundle.ChainID,
		},
	}, nil
}

// Config returns the configuration for this PluginBackend.
func (b *PluginBackend) readConfig(ctx context.Context, s logical.Storage) (*ConfigJSON, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf("the plugin has not been configured yet")
	}

	var result ConfigJSON
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, fmt.Errorf("error reading configuration: %s", err)
		}
	}

	return &result, nil
}

func (b *PluginBackend) configured(ctx context.Context, req *logical.Request) (*ConfigJSON, error) {
	config, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if validConnection, err := b.validIPConstraints(config, req); !validConnection {
		return nil, err
	}

	return config, nil
}

func (b *PluginBackend) validIPConstraints(config *ConfigJSON, req *logical.Request) (bool, error) {
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
