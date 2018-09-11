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

// swagger:parameters pathAddressesRead pathAccountBalanceReadByAddress
type AddressRequest struct {
	// The address to lookup
	//
	// in: path
	// required: true
	Address string `json:"address"`
}

// swagger:parameters pathAddressesList pathAccountsList
type ListRequest struct {
	// So that we can get the list from Vault.  Do not change this.
	//
	// in: query
	// required: true
	// default: true
	List bool `json:"list"`
}

// swagger:parameters pathAddressesVerify
type AddressVerifyRequest struct {
	// The address to lookup
	//
	// in: path
	// required: true
	Address string `json:"address"`
	// in: body
	// required: true
	// schema:
	//	 type: string
	Data struct {
		Data      string `json:"data"`
		Signature string `json:"signature"`
	} `json:"data"`
}

// swagger:parameters pathCreateConfig
type ConfigRequest struct {
	// The conversion inputs
	//
	// in: body
	// required: true
	Data struct {
		RpcUrl        string   `json:"rpc_url,omitempty"`
		ChainId       string   `json:"chain_id,omitempty"`
		BoundCidrList []string `json:"bound_cidr_list,omitempty"`
		ApiKey        string   `json:"api_key,omitempty"`
	} `json:"data"`
}

// swagger:parameters pathConvertWrite
type ConversionRequest struct {
	// The conversion inputs
	//
	// in: body
	// required: true
	// schema:
	//	 type: string
	Data struct {
		AmountIn   string `json:"amount"`
		UnitFromIn string `json:"unit_from"`
		UnitToIn   string `json:"unit_to"`
	} `json:"data"`
}

// BaseResponse stores the names of the account to allow reverse lookup by address
type BaseResponse struct {
	RequestId     string   `json:"request_id"`
	LeaseId       string   `json:"lease_id"`
	Renewable     bool     `json:"renewable"`
	LeaseDuration int      `json:"lease_duration"`
	WrapInfo      string   `json:"wrap_info"`
	Warnings      []string `json:"warnings"`
	Auth          string   `json:"auths"`
}

// Returns a list of keys
// swagger:model KeyListResponse
type KeyListResponse struct {
	BaseResponse
	Data struct {
		Keys []string `json:"keys"`
	} `json:"data"`
}

// swagger:model AddressBalanceResponse
type AddressBalanceResponse struct {
	BaseResponse
	Data struct {
		Address      string `json:"address"`
		Balance      string `json:"balance"`
		BalanceInUsd bool   `json:"balance_in_usd"`
	}
}

// swagger:model AccountNamesResponse
type AccountNamesResponse struct {
	BaseResponse
	Data struct {
		Names []string `json:"names"`
	} `json:"data"`
}

// swagger:model AddressesVerifiedResponse
type AddressesVerifiedResponse struct {
	BaseResponse
	Data struct {
		Address   string `json:"address"`
		Signature string `json:"signature"`
		Verified  bool   `json:"verified"`
	} `json:"data"`
}

// swagger:model ConfigResponse
type ConfigResponse struct {
	BaseResponse
	Data struct {
		ApiKey        string   `json:"api_key"`
		BoundCidrList []string `json:"bound_cidr_list"`
		ChainId       string   `json:"chain_id"`
		RpcUrl        string   `json:"rpc_url"`
	} `json:"data"`
}

// swagger:model ConversionResponse
type ConversionResponse struct {
	BaseResponse
	Data struct {
		AmountFrom string `json:"amount_from"`
		AmountTo   string `json:"amount_to"`
		UnitFrom   string `json:"unit_from"`
		UnitTo     string `json:"unit_to"`
	} `json:"data"`
}
