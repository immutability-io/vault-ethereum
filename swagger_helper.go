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

// swagger:parameters pathCreateConfig pathUpdateConfig pathReadConfig pathConvertWrite pathAddressesRead pathAddressesList pathAddressesVerify pathAccountBalanceReadByAddress pathAccountsList pathAccountsDelete pathAccountsRead pathAccountsCreate pathAccountUpdate pathVerify pathDebit pathContractsList pathTransfer pathSign pathCreateContract pathReadContract pathContractsDelete pathBlockRead pathBlockTransactionsList pathExportCreate pathImportCreate pathNamesList pathNamesRead pathNamesVerify pathTransactionRead
type MountPathParam struct {
	// The endpoint configured for the plugin mount
	//
	// in: path
	// required: true
	MountPath string `json:"mount-path"`
}

// swagger:parameters pathAccountsDelete pathAccountsRead pathAccountsCreate pathAccountUpdate pathVerify pathDebit pathContractsList pathTransfer pathSign pathCreateContract pathReadContract pathContractsDelete pathExportCreate pathImportCreate pathNamesRead pathNamesVerify
type AccountNameParam struct {
	// The account name
	//
	// in: path
	// required: true
	AccountName               string `json:"name"`
}

// swagger:parameters pathAccountsCreate pathAccountUpdate
type AccountRequest struct {
	// The account to modify
	//
	// in: body
	// required: true
	// schema:
	//	 type: string
	Data struct {
		SpendingLimitTx    string `json:"spending_limit_tx,omitempty"`
		SpendingLimitTotal string `json:"spending_limit_total,omitempty"`
		Whitelist          string `json:"whitelist,omitempty"`
		Blacklist          string `json:"blacklist,omitempty"`
	}
}

// swagger:parameters pathAddressesRead pathAccountBalanceReadByAddress pathAddressesVerify
type AddressParam struct {
	// The address to lookup
	//
	// in: path
	// required: true
	Address string `json:"address"`
}

// swagger:parameters pathBlockRead pathBlockTransactionsList
type BlockNumberParam struct {
	// The block number
	//
	// in: path
	// required: true
	BlockNumber               string `json:"block-number"`
}


// swagger:parameters pathCreateContract pathReadContract pathContractsDelete
type ContractNameParam struct {
	// The contract name
	//
	// in: path
	// required: true
	ContractName               string `json:"contract-name"`
}

// swagger:parameters pathExportCreate
type ExportRequest struct {
	// The path to export to
	//
	// in: body
	// required: true
	// schema:
	//	 type: string
	Data struct {
		ExportPath               string `json:"path"`
	} `json:"data"`
}

// swagger:parameters pathImportCreate
type ImportRequest struct {
	// The data to import from
	//
	// in: body
	// required: true
	// schema:
	//	 type: string
	Data struct {
		ImportPath               string `json:"path"`
		ImportPassphrase               string `json:"passphrase"`
	} `json:"data"`
}

// swagger:parameters pathAddressesList pathAccountsList pathContractsList pathNamesList
type ListRequest struct {
	// So that we can get the list from Vault.  Do not change this.
	//
	// in: query
	// required: true
	// default: true
	List bool `json:"list"`
}

// swagger:parameters pathCreateConfig pathUpdateConfig
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


// swagger:parameters pathCreateContract
type ContractRequest struct {
	// The contract inputs
	//
	// in: body
	// required: true
	// schema:
	//	 type: string
	Data struct {
		TransactionData   string `json:"transaction_data"`
		Amount string `json:"amount"`
		Nonce   string `json:"nonce,omitempty"`
		GasPrice string `json:"gas_price"`
		GasLimit string `json:"gas_limit"`
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

// swagger:parameters pathCreateConfig pathDebit
type DebitRequest struct {
	// The debit inputs
	//
	// in: body
	// required: true
	Data struct {
		AddressTo        string   `json:"address_to"`
		Amount       string   `json:"amount"`
		GasPrice []string `json:"gas_price"`
		GasLimit        string   `json:"gas_limit"`
	} `json:"data"`
}

// swagger:parameters pathAddressesVerify pathTransactionRead
type TransactionHashParam struct {
	// in: path
	// required: true
	TransactionHash string `json:"transaction_hash"`
}

// swagger:parameters pathAddressesVerify pathNamesVerify
type VerifyRequest struct {
	// in: body
	// required: true
	// schema:
	//	 type: string
	Data struct {
		Data      string `json:"data"`
		Signature string `json:"signature"`
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

// swagger:model AddressListResponse
type AddressListResponse struct {
	BaseResponse
	Data struct {
		Address      []string `json:"address"`
	}
}

// swagger:model AccountResponse
type AccountReponse struct {
	BaseResponse
	Data struct {
		Address      string `json:"address"`
		Balance      string `json:"balance"`
		BalanceInUsd bool   `json:"balance_in_usd"`
		Blacklist string `json:"blacklist"`
		SpendingLimitTotal string `json:"spending_limit_total"`
		SpendingLimitTx string `json:"spending_limit_tx"`
		TotalSpend string `json:"total_spend"`
		Whitelist string `json:"whitelist"`
	}
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

// swagger:model BlockResponse
type BlockResponse struct {
	BaseResponse
	Data struct {
		Block string `json:"block"`
		BlockHash string `json:"block_hash"`
		Difficulty int `json:"difficulty"`
		Time string `json:"time"`
		TransactionCount string `json:"transaction_count"`
	} `json:"data"`
}


// swagger:model BlockTransactionsResponse
type BlockTransactionsResponse struct {
	BaseResponse
	Data []struct {
		Block struct {
			AddressTo string `json:"address_to"`
		} `json:"block"`
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

// swagger:model ContractResponse
type ContractResponse struct {
	BaseResponse
	Data struct {
		Address string `json:"transaction_hash,omitempty"`
		TransactionHash string `json:"transaction_hash"`

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

// swagger:model DebitResponse
type DebitResponse struct {
	BaseResponse
	Data struct {
		Amount       string   `json:"amount"`
		FromAddress        string   `json:"from_address"`
		GasLimit        string   `json:"gas_limit"`
		GasPrice string `json:"gas_price"`
		Balance string `json:"balance"`
		ToAddress string `json:"to_address"`
		TotalSpend string `json:"total_spend"`
		TransactionHash string `json:"transaction_hash"`

	} `json:"data"`
}

// swagger:model ExportResponse
type ExportResponse struct {
	BaseResponse
	Data struct {
		Passphrase       string   `json:"passphrase"`
		Path        string   `json:"path"`

	} `json:"data"`
}

// swagger:model SignedResponse
type SignedResponse struct {
	BaseResponse
	Data struct {
		Address   string `json:"address"`
		Signature string `json:"signature"`
	} `json:"data"`
}

// swagger:model TransactionResponse
type TransactionResponse struct {
	BaseResponse
	Data struct {
		AddressFrom   string `json:"address_from"`
		AddressTo   string `json:"address_to"`
		Gas   string `json:"gas"`
		GasPrice   string `json:"gas_price"`
		Nonce   int `json:"nonce"`
		Pending   bool `json:"pending"`
		ReceiptStatus   string `json:"receipt_status"`
		TransactionHash   string `json:"transaction_hash"`
		Value   string `json:"value"`
	} `json:"data"`
}

// swagger:model VerifiedResponse
type VerifiedResponse struct {
	BaseResponse
	Data struct {
		Address   string `json:"address"`
		Signature string `json:"signature"`
		Verified  bool   `json:"verified"`
	} `json:"data"`
}
