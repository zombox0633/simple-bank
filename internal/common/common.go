package common

import "errors"

const (
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	CurrencyTHB = "THB"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("resource already exists")
	ErrInvalidReference    = errors.New("referenced resource not found")
	ErrUnsupportedCurrency = errors.New("unsupported currency")
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case CurrencyUSD, CurrencyEUR, CurrencyTHB:
		return true
	default:
		return false
	}
}
