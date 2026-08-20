package transfer

import "github.com/govalues/decimal"

type createTransferRequest struct {
	FromAccountID int64           `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64           `json:"to_account_id" binding:"required,min=1,nefield=FromAccountID"`
	Amount        decimal.Decimal `json:"amount" binding:"money"`
	Currency      string          `json:"currency" binding:"required,currency"`
}
