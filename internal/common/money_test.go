package common

import (
	"testing"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"
)

func TestIsValidMoneyAmount(t *testing.T) {
	tests := []struct {
		amount string
		valid  bool
	}{
		{"10", true},
		{"10.1", true},
		{"10.1234", true},
		{"10.12340", true},
		{"99999999999999.9999", true},
		{"10.12345", false},
		{"100000000000000", false},
		{"0", false},
		{"-10", false},
	}

	for _, test := range tests {
		t.Run(test.amount, func(t *testing.T) {
			amount := decimal.MustParse(test.amount)
			require.Equal(t, test.valid, IsValidMoneyAmount(amount))
		})
	}
}
