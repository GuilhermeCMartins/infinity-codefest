package utils

import (
	"fmt"
	"math/big"
)

func Float64ToBigFloat(f float64) *big.Float {
	return new(big.Float).SetFloat64(f)
}

func BigFloatToFloat64(f *big.Float) (float64, error) {
	result, accuracy := f.Float64()
	if accuracy == big.Exact {
		return result, nil
	}
	return 0, fmt.Errorf("precision loss during conversion")
}
