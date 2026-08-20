package api

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"

	"simplebank/internal/common"
)

func RegisterValidators() {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	_ = validate.RegisterValidation("currency", validCurrency)
	_ = validate.RegisterValidation("money", validMoney)
}

func validCurrency(field validator.FieldLevel) bool {
	currency, ok := field.Field().Interface().(string)
	return ok && common.IsSupportedCurrency(currency)
}

func validMoney(field validator.FieldLevel) bool {
	amount, ok := field.Field().Interface().(decimal.Decimal)
	return ok && common.IsValidMoneyAmount(amount)
}
