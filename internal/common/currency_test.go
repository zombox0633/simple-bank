package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSupportedCurrency(t *testing.T) {
	tests := map[string]bool{
		CurrencyUSD: true,
		CurrencyEUR: true,
		CurrencyTHB: true,
		"BTC":       false,
	}

	for currency, supported := range tests {
		t.Run(currency, func(t *testing.T) {
			require.Equal(t, supported, IsSupportedCurrency(currency))
		})
	}
}
