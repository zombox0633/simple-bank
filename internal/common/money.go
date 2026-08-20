package common

import "github.com/govalues/decimal"

const MoneyScale = 4

var MaxMoneyAmount = decimal.MustParse("99999999999999.9999")

// IsValidMoneyAmount accepts positive amounts up to MaxMoneyAmount with at most MoneyScale decimal places.
// Whole numbers such as 10 are valid; NUMERIC(18,4) stores them as 10.0000.
func IsValidMoneyAmount(amount decimal.Decimal) bool {
	return amount.IsPos() &&
		amount.Cmp(MaxMoneyAmount) <= 0 &&
		amount.MinScale() <= MoneyScale
}
