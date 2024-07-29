package utils

import (
	"fmt"
)

type Currency string

const (
	IC  Currency = "ic"
	USD Currency = "usd"
	BR  Currency = "br"
	ETH Currency = "eth"
	BTC Currency = "btc"
)

// Define a map with conversion rates relative to USD
var conversionRates = map[Currency]float64{
	USD: 1,
	BR:  0.2,
	ETH: 3000,
	BTC: 60000,
	IC:  1000000,
}

func ConvertCurrency(amount float64, senderCurrency Currency, transactionCurrency Currency, receiverCurrency Currency) (float64, float64, error) {
	convertToUSD := func(amount float64, fromCurrency Currency) (float64, error) {
		fromRate, exists := conversionRates[fromCurrency]
		if !exists {
			return 0, fmt.Errorf("Conversion tax not found %s", fromCurrency)
		}
		return amount * fromRate, nil
	}

	convertFromUSD := func(amount float64, toCurrency Currency) (float64, error) {
		toRate, exists := conversionRates[toCurrency]
		if !exists {
			return 0, fmt.Errorf("Conversion tax not found %s", toCurrency)
		}
		return amount / toRate, nil
	}

	amountInUSD, err := convertToUSD(amount, transactionCurrency)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid conversion from %s to USD: %v", transactionCurrency, err)
	}

	amountInSenderCurrency, err := convertFromUSD(amountInUSD, senderCurrency)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid conversion from USD to %s: %v", senderCurrency, err)
	}

	amountInReceiverCurrency, err := convertFromUSD(amountInUSD, receiverCurrency)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid conversion from USD to %s: %v", receiverCurrency, err)
	}

	return amountInReceiverCurrency, amountInSenderCurrency, nil
}
