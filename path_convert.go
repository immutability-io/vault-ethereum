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
	"log"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
	"github.com/shopspring/decimal"
)

const (
	// WEI is "wei"
	WEI string = "wei"
	// KWEI is "ada", "kwei", "kilowei", "femtoether"
	KWEI string = "kwei"
	// MWEI is "babbage", "mwei", "megawei", "picoether"
	MWEI string = "mwei"
	// GWEI is "shannon", "gwei", "gigawei", "nanoether", "nano"
	GWEI string = "gwei"
	// MICRO is "szazbo", "micro", "microether"
	MICRO string = "micro"
	// MILLI is "finney", "milli", "milliether"
	MILLI string = "milli"
	// ETH is "ether", "eth"
	ETH string = "ether"
	// KILO is "einstein", "kilo", "kiloether", "kether", "grand"
	KILO string = "kilo"
	// MEGA is "mega", "megaether", "mether"
	MEGA string = "mega"
	// GIGA is giga", "gigaether", "gether"
	GIGA string = "giga"
	// TERA is "tera", "teraether", "tether"
	TERA string = "tera"
	// USD is "usd"
	USD string = "usd"
	// EUR is "eur"
	EUR string = "eur"
)

func convertPaths(b *EthereumBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:      "convert",
			HelpSynopsis: "Convert any Ethereum unit to another.",
			HelpDescription: `

			Convert any Ethereum unit to another.
`,
			Fields: map[string]*framework.FieldSchema{
				"unit_from": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The Ethereum unit to convert from.",
				},
				"unit_to": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "The Ethereum unit to convert to.",
				},
				"amount": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Amount to convert.",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathConvertWrite,
				logical.UpdateOperation: b.pathConvertWrite,
			},
		},
	}
}

// ValidUnit returns a normalized metric unit or an error
func ValidUnit(unit string) (string, error) {
	switch strings.ToLower(unit) {
	case "wei":
		return WEI, nil
	case "ada", "kwei", "kilowei", "femtoether":
		return KWEI, nil
	case "babbage", "mwei", "megawei", "picoether":
		return MWEI, nil
	case "shannon", "gwei", "gigawei", "nanoether", "nano":
		return GWEI, nil
	case "szazbo", "micro", "microether":
		return MICRO, nil
	case "finney", "milli", "milliether":
		return MILLI, nil
	case "ether", "eth":
		return ETH, nil
	case "einstein", "kilo", "kiloether", "kether", "grand":
		return KILO, nil
	case "mega", "megaether", "mether":
		return MEGA, nil
	case "giga", "gigaether", "gether":
		return GIGA, nil
	case "tera", "teraether", "tether":
		return TERA, nil
	case "usd", "USD":
		return USD, nil
	}
	return "", fmt.Errorf("Unknown unit %s", unit)
}

// ToWeiMultiplier returns the multipler to convert a unit to wei
func ToWeiMultiplier(normalizeUnit string) decimal.Decimal {
	var multiplier decimal.Decimal
	switch normalizeUnit {
	case WEI:
		multiplier, _ = decimal.NewFromString("1")
	case KWEI:
		multiplier, _ = decimal.NewFromString("1000")
	case MWEI:
		multiplier, _ = decimal.NewFromString("1000000")
	case GWEI:
		multiplier, _ = decimal.NewFromString("1000000000")
	case MICRO:
		multiplier, _ = decimal.NewFromString("1000000000000")
	case MILLI:
		multiplier, _ = decimal.NewFromString("1000000000000000")
	case ETH:
		multiplier, _ = decimal.NewFromString("1000000000000000000")
	case KILO:
		multiplier, _ = decimal.NewFromString("1000000000000000000000")
	case MEGA:
		multiplier, _ = decimal.NewFromString("1000000000000000000000000")
	case GIGA:
		multiplier, _ = decimal.NewFromString("1000000000000000000000000000")
	case TERA:
		multiplier, _ = decimal.NewFromString("1000000000000000000000000000000")
	}
	return multiplier
}

// FromWeiMultiplier returns the multipler to convert a unit to wei
func FromWeiMultiplier(normalizeUnit string) decimal.Decimal {
	var multiplier decimal.Decimal
	switch normalizeUnit {
	case WEI:
		multiplier, _ = decimal.NewFromString("1")
	case KWEI:
		multiplier, _ = decimal.NewFromString("0.001")
	case MWEI:
		multiplier, _ = decimal.NewFromString("0.000001")
	case GWEI:
		multiplier, _ = decimal.NewFromString("0.000000001")
	case MICRO:
		multiplier, _ = decimal.NewFromString("0.000000000001")
	case MILLI:
		multiplier, _ = decimal.NewFromString("0.000000000000001")
	case ETH:
		multiplier, _ = decimal.NewFromString("0.000000000000000001")
	case KILO:
		multiplier, _ = decimal.NewFromString("0.000000000000000000001")
	case MEGA:
		multiplier, _ = decimal.NewFromString("0.000000000000000000000001")
	case GIGA:
		multiplier, _ = decimal.NewFromString("0.000000000000000000000000001")
	case TERA:
		multiplier, _ = decimal.NewFromString("0.000000000000000000000000000001")
	}
	return multiplier
}

// ConvertToWei will convert any valid ethereum unit to wei
func ConvertToWei(normalizeUnit string, amount decimal.Decimal) decimal.Decimal {
	var result decimal.Decimal
	result = amount.Mul(ToWeiMultiplier(normalizeUnit))
	return result
}

// ConvertFromWei will convert any valid ethereum unit to wei
func ConvertFromWei(normalizeUnit string, amount decimal.Decimal) decimal.Decimal {
	var result decimal.Decimal
	result = amount.Mul(FromWeiMultiplier(normalizeUnit))
	return result
}

// ConvertToUSD uses Coinmarketcap to estimate value of ETH in USD
func ConvertToUSD(amountInWei, apiKey string) (decimal.Decimal, error) {
	zero, _ := decimal.NewFromString("0")
	if apiKey == "" {
		return zero, nil
	}

	client := cmc.NewClient(&cmc.Config{
		ProAPIKey: apiKey,
	})

	listings, err := client.Cryptocurrency.LatestListings(&cmc.ListingOptions{
		Limit: 10,
	})
	if err != nil {
		log.Fatal(err)
	}
	var ethListing *cmc.Listing
	for _, listing := range listings {
		if listing.Symbol == "ETH" {
			ethListing = listing
			break
		}
	}
	if ethListing == nil {
		return zero, nil
	}
	balanceInWei, err := decimal.NewFromString(amountInWei)
	if err != nil {
		return zero, err
	}
	quote, ok := ethListing.Quote["USD"]
	if ok {
		price := decimal.NewFromFloat(quote.Price)
		balanceInETH := ConvertFromWei(ETH, balanceInWei)
		exchangeValue := price.Mul(balanceInETH)
		return exchangeValue, nil
	}
	return zero, nil

}

func (b *EthereumBackend) pathConvertWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	usdFrom := false
	usdTo := false
	exchangeValue, _ := decimal.NewFromString("0")
	amountFromInUnits, _ := decimal.NewFromString("0")
	oneETH, _ := decimal.NewFromString("1")
	unitFrom, err := ValidUnit(data.Get("unit_from").(string))
	if err != nil {
		return nil, err
	}
	amountFrom := data.Get("amount").(string)
	amount, err := decimal.NewFromString(amountFrom)
	if err != nil || amount.IsNegative() {
		return nil, fmt.Errorf("amount is either not a number or is negative")
	}

	unitTo, err := ValidUnit(data.Get("unit_to").(string))
	if err != nil {
		return nil, err
	}
	if unitFrom == unitTo {
		return nil, fmt.Errorf("Conversion from %s to %s makes no sense", unitFrom, unitTo)
	}
	if unitFrom == USD || unitTo == USD {
		oneETHInWei := ConvertToWei(ETH, oneETH)
		exchangeValue, err = ConvertToUSD(oneETHInWei.String(), config.CoinMarketCapAPIKey)
		if err != nil {
			return nil, err
		}
	}
	if unitFrom == USD {
		usdFrom = true
		unitFrom = ETH
		amount = amount.Div(exchangeValue)
		if err != nil {
			return nil, err
		}
	}
	amountFromInWei := ConvertToWei(unitFrom, amount)
	if unitTo == USD {
		usdTo = true
		ethInWei := ConvertFromWei(ETH, amountFromInWei)
		amountFromInUnits = ethInWei.Mul(exchangeValue)
	} else {
		amountFromInUnits = ConvertFromWei(unitTo, amountFromInWei)
	}

	if usdFrom {
		unitFrom = data.Get("unit_from").(string)
		amount, _ = decimal.NewFromString(data.Get("amount").(string))
	} else if usdTo {
		unitTo = data.Get("unit_to").(string)
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"unit_from":   unitFrom,
			"amount_from": amount,
			"unit_to":     unitTo,
			"amount_to":   amountFromInUnits.String(),
		},
	}, nil
}
