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

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// Symbol is the lowercase crypto token symbol
	Symbol string = "eth"
)

// Factory returns the backend
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// FactoryType returns the factory
func FactoryType(backendType logical.BackendType) logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		b, err := Backend(conf)
		if err != nil {
			return nil, err
		}
		b.BackendType = backendType
		if err = b.Setup(ctx, conf); err != nil {
			return nil, err
		}
		return b, nil
	}
}

// Backend returns the backend
func Backend(conf *logical.BackendConfig) (*PluginBackend, error) {
	var b PluginBackend
	b.Backend = &framework.Backend{
		Help: "",
		Paths: framework.PathAppend(
			configPaths(&b),
			accountPaths(&b),
			convertPaths(&b),
			erc20Paths(&b),
		),
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"convert",
				"test",
			},
			SealWrapStorage: []string{
				"accounts/",
			},
		},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}
	return &b, nil
}

// PluginBackend implements the Backend for this plugin
type PluginBackend struct {
	*framework.Backend
}

// QualifiedPath prepends the token symbol to the path
func QualifiedPath(subpath string) string {
	return subpath
}

// ContractPath prepends the token symbol to the path
func ContractPath(contract, method string) string {
	return fmt.Sprintf("%s/%s/%s", QualifiedPath("accounts/"+framework.GenericNameRegex("name")), contract, method)
}

// SealWrappedPaths returns the paths that are seal wrapped
func SealWrappedPaths(b *PluginBackend) []string {
	return []string{
		QualifiedPath("accounts/"),
	}
}
