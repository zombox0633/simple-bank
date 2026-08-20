package common

import (
	"errors"

	"github.com/shopspring/decimal"
)

const (
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	CurrencyTHB = "THB"
	MoneyScale  = 4
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("resource already exists")
	ErrInvalidReference    = errors.New("referenced resource not found")
	ErrUnsupportedCurrency = errors.New("unsupported currency")
	ErrCurrencyMismatch    = errors.New("account currency does not match transfer currency")
	ErrInvalidAmount       = errors.New("amount must be greater than zero and have at most 4 decimal places")
	ErrSameAccount         = errors.New("source and destination accounts must be different")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case CurrencyUSD, CurrencyEUR, CurrencyTHB:
		return true
	default:
		return false
	}
}

func IsValidMoneyAmount(amount decimal.Decimal) bool {
	return amount.GreaterThan(decimal.Zero) && amount.Equal(amount.Round(MoneyScale))
}
