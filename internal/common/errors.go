package common

import "errors"

var (
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
