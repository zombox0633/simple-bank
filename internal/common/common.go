package common

import (
	"errors"

	"github.com/govalues/decimal"
)

const (
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	CurrencyTHB = "THB"
	MoneyScale  = 4
)

var (
	MaxMoneyAmount          = decimal.MustParse("99999999999999.9999")
	ErrNotFound             = errors.New("resource not found")
	ErrConflict             = errors.New("resource already exists")
	ErrInvalidReference     = errors.New("referenced resource not found")
	ErrUnsupportedCurrency  = errors.New("unsupported currency")
	ErrCurrencyMismatch     = errors.New("account currency does not match transfer currency")
	ErrInvalidAmount        = errors.New("amount must be greater than zero, not exceed 99999999999999.9999, and have at most 4 decimal places")
	ErrSameAccount          = errors.New("source and destination accounts must be different")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrBalanceLimitExceeded = errors.New("destination account balance limit exceeded")
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case CurrencyUSD, CurrencyEUR, CurrencyTHB:
		return true
	default:
		return false
	}
}

// IsValidMoneyAmount accepts positive amounts up to MaxMoneyAmount with at most MoneyScale decimal places.
// Whole numbers such as 10 are valid; NUMERIC(18,4) stores them as 10.0000.
func IsValidMoneyAmount(amount decimal.Decimal) bool {
	return amount.IsPos() &&
		amount.Cmp(MaxMoneyAmount) <= 0 &&
		amount.MinScale() <= MoneyScale
}
