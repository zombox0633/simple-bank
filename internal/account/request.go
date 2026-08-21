package account

import (
	"math"

	z "github.com/Oudwins/zog"

	"simplebank/internal/httpapi"
)

type createAccountRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
}

var createAccountSchema = z.Struct(z.Shape{
	"owner": z.String().Trim().Transform(httpapi.Lowercase).Min(1).Required(),
	"currency": z.String().Trim().Transform(httpapi.Uppercase).TestFunc(
		httpapi.IsSupportedCurrency,
		z.IssueCode(httpapi.IssueUnsupportedCurrency),
	).Required(z.IssueCode(httpapi.IssueUnsupportedCurrency)),
})

type getAccountRequest struct {
	ID int64 `uri:"id" zog:"id"`
}

var getAccountSchema = z.Struct(z.Shape{
	"ID": z.Int64().GT(0).Required(),
})

type listAccountsRequest struct {
	Owner    string `form:"owner" zog:"owner"`
	PageID   int32  `form:"page_id" zog:"page_id"`
	PageSize int32  `form:"page_size" zog:"page_size"`
}

var listAccountsSchema = z.Struct(z.Shape{
	"owner":    z.String().Trim().Transform(httpapi.Lowercase).Min(1).Required(),
	"pageID":   z.Int32().GTE(1).Required(),
	"pageSize": z.Int32().GTE(1).LTE(100).Required(),
}).TestFunc(func(value any, _ z.Ctx) bool {
	request := value.(*listAccountsRequest)
	if request.PageID <= 0 || request.PageSize <= 0 {
		return true // field schemas report these errors
	}
	offset := int64(request.PageID-1) * int64(request.PageSize)
	return offset <= math.MaxInt32
}, z.IssuePath([]string{"page_id"}))
