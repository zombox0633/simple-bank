package common

import "errors"

var (
	ErrNotFound             = errors.New("resource not found")
	ErrConflict             = errors.New("resource already exists")
	ErrInvalidReference     = errors.New("referenced resource not found")
	ErrInvalidUsername      = errors.New("username must contain 3 to 32 characters, start with a letter, contain only letters, numbers, or underscores, and not end with an underscore")
	ErrInvalidFullName      = errors.New("full name must contain 1 to 100 letters and may include spaces, apostrophes, hyphens, or periods")
	ErrInvalidEmail         = errors.New("email address is invalid or exceeds 254 bytes")
	ErrUnsupportedCurrency  = errors.New("unsupported currency")
	ErrCurrencyMismatch     = errors.New("account currency does not match transfer currency")
	ErrInvalidAmount        = errors.New("amount must be greater than zero, not exceed 99999999999999.9999, and have at most 4 decimal places")
	ErrInvalidPassword      = errors.New("password must contain 8 to 64 characters with at least one letter, one number, and one special character")
	ErrSameAccount          = errors.New("source and destination accounts must be different")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrBalanceLimitExceeded = errors.New("destination account balance limit exceeded")
)
