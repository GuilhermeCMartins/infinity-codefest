package utils

import "fmt"

type Currency string

const (
	IC  Currency = "IC"
	USD Currency = "USD"
	BRL Currency = "BRL"
	ETH Currency = "ETH"
	BTC Currency = "BTC"
)

func ConvertCurrency(amount float64, requestCurrency Currency, transactionCurrency Currency, receiverCurrency Currency) (float64, float64, error) {
	conversionRates := map[string]float64{
		"USD": 1,
		"BRL": 0.2,
		"ETH": 3000,
		"BTC": 60000,
		"IC":  1000000,
	}

	convert := func(amount float64, fromCurrency Currency, toCurrency Currency) (float64, error) {
		fromRate, fromExists := conversionRates[string(fromCurrency)]
		toRate, toExists := conversionRates[string(toCurrency)]

		if !fromExists || !toExists {
			return 0, fmt.Errorf("Invalid currency: %s ou %s", fromCurrency, toCurrency)
		}

		amountInUSD := amount / fromRate
		convertedAmount := amountInUSD * toRate

		return convertedAmount, nil
	}

	var err error
	var intermediateAmount float64

	if requestCurrency != transactionCurrency {
		intermediateAmount, err = convert(amount, requestCurrency, transactionCurrency)
		if err != nil {
			return 0, 0, fmt.Errorf("erro na conversão de %s para %s: %v", requestCurrency, transactionCurrency, err)
		}
	} else {
		intermediateAmount = amount
	}

	finalAmount := intermediateAmount
	if transactionCurrency != receiverCurrency {
		finalAmount, err = convert(intermediateAmount, transactionCurrency, receiverCurrency)
		if err != nil {
			return 0, 0, fmt.Errorf("erro na conversão de %s para %s: %v", transactionCurrency, receiverCurrency, err)
		}
	}

	return intermediateAmount, finalAmount, nil
}
