package transfer

import (
	z "github.com/Oudwins/zog"
	"github.com/govalues/decimal"

	"simplebank/internal/common"
	"simplebank/internal/httpapi"
)

type createTransferRequest struct {
	FromAccountID int64           `json:"from_account_id"`
	ToAccountID   int64           `json:"to_account_id"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
}

var moneySchema = z.CustomFunc(
	func(amount *decimal.Decimal, _ z.Ctx) bool {
		return common.IsValidMoneyAmount(*amount)
	},
	z.IssueCode(httpapi.IssueInvalidAmount),
)

var createTransferSchema = z.Struct(z.Shape{
	"fromAccountID": z.Int64().GT(0).Required(),
	"toAccountID":   z.Int64().GT(0).Required(),
	"amount":        moneySchema,
	"currency": z.String().Trim().Transform(httpapi.Uppercase).TestFunc(
		httpapi.IsSupportedCurrency,
		z.IssueCode(httpapi.IssueUnsupportedCurrency),
	).Required(z.IssueCode(httpapi.IssueUnsupportedCurrency)),
}).TestFunc(
	func(value any, _ z.Ctx) bool {
		request := value.(*createTransferRequest)
		return request.FromAccountID <= 0 ||
			request.ToAccountID <= 0 ||
			request.FromAccountID != request.ToAccountID
	},
	z.IssueCode(httpapi.IssueSameAccount),
	z.IssuePath([]string{"to_account_id"}),
)
