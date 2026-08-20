package common

const (
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	CurrencyTHB = "THB"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case CurrencyUSD, CurrencyEUR, CurrencyTHB:
		return true
	default:
		return false
	}
}
