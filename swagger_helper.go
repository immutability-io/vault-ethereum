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


// swagger:parameters pathAddressesRead
type AddressParam struct {
	// The address to lookup
	//
	// in: path
	// required: true
	Address string `json:"address"`
}


// swagger:parameters pathAddressesList
type AddressListParams struct {
	// So that we can get the list from Vault.  Do not change this.
	//
	// in: query
	// required: true
	// default: true
	List bool `json:"list"`
}

// swagger:parameters pathAddressesVerify
type AddressVerifyParams struct {
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
		Data string `json:"data"`
		Signature string `json:"signature"`
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
	Data struct  {
		AmountIn string `json:"amount"`
		UnitFromIn string `json:"unit_from"`
		UnitToIn string `json:"unit_to"`
	} `json:"data"`
}

// BaseStruct stores the names of the account to allow reverse lookup by address
// swagger:model baseStruct
type BaseStruct struct {
	RequestId string `json:"request_id"`
	LeaseId string `json:"lease_id"`
	Renewable bool `json:"renewable"`
	LeaseDuration int `json:"lease_duration"`
	WrapInfo string `json:"wrap_info"`
	Warnings []string `json:"warnings"`
	Auth string `json:"auths"`
}

// Addresses stores the names of the account to allow reverse lookup by address
// swagger:model AddressesResponse
type AddressesResponse struct {
	BaseStruct
	Data struct {
		Keys []string `json:"keys"`
	}  `json:"data"`
}

// AccountNamesResponse stores the list of addresses
// swagger:model AccountNamesResponse
type AccountNamesResponse struct {
	BaseStruct
	Data struct {
		Names []string `json:"names"`
	} `json:"data"`
}

// Addresses stores the status of an Address verification response
// swagger:model AddressesVerifiedResponse
type AddressesVerifiedResponse struct {
	BaseStruct
	Data struct {
		Address string `json:"address"`
		Signature string `json:"signature"`
		Verified bool `json:"verified"`
	} `json:"data"`
}

// Addresses stores the calculated conversion data
// swagger:model ConversionResponse
type SwaggerConvertStruct struct {
	BaseStruct
	Data struct {
		AmountFrom string `json:"amount_from"`
		AmountTo string `json:"amount_to"`
		UnitFrom string `json:"unit_from"`
		UnitTo string `json:"unit_to"`
	} `json:"data"`
}