package utils

import "fmt"

type Currency string

const (
	IC  Currency = "ic"
	USD Currency = "usd"
	BRL Currency = "brl"
	ETH Currency = "eth"
	BTC Currency = "btc"
)

func ConvertCurrency(amount float64, requestCurrency Currency, transactionCurrency Currency, receiverCurrency Currency) (float64, float64, error) {
	conversionRates := map[string]float64{
		"usd": 1,
		"brl": 0.2,
		"eth": 3000,
		"btc": 60000,
		"ic":  1000000,
	}

	print(requestCurrency)
	print(transactionCurrency)
	print(receiverCurrency)

	convert := func(amount float64, fromCurrency Currency, toCurrency Currency) (float64, error) {
		print(fromCurrency)
		print(toCurrency)
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
