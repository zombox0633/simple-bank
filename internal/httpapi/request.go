package httpapi

import (
	"net/http"
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/gin-gonic/gin"

	"simplebank/internal/common"
)

const (
	IssueInvalidUsername     = "invalid_username"
	IssueInvalidFullName     = "invalid_full_name"
	IssueInvalidEmail        = "invalid_email"
	IssueInvalidPassword     = "invalid_password"
	IssueUnsupportedCurrency = "unsupported_currency"
	IssueInvalidAmount       = "invalid_amount"
	IssueSameAccount         = "same_account"
)

func BindJSON(ctx *gin.Context, schema *z.StructSchema, destination any) bool {
	return bind(ctx, ctx.ShouldBindJSON, schema, destination)
}

func BindQuery(ctx *gin.Context, schema *z.StructSchema, destination any) bool {
	return bind(ctx, ctx.ShouldBindQuery, schema, destination)
}

func BindURI(ctx *gin.Context, schema *z.StructSchema, destination any) bool {
	return bind(ctx, ctx.ShouldBindUri, schema, destination)
}

func bind(
	ctx *gin.Context,
	binder func(any) error,
	schema *z.StructSchema,
	destination any,
) bool {
	if err := binder(destination); err != nil {
		writeError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return false
	}
	issues := schema.Validate(destination)
	if len(issues) == 0 {
		return true
	}
	WriteRequestIssues(ctx, issues)
	return false
}

func Lowercase(value *string, _ z.Ctx) error {
	*value = strings.ToLower(*value)
	return nil
}

func Uppercase(value *string, _ z.Ctx) error {
	*value = strings.ToUpper(*value)
	return nil
}

func CollapseWhitespace(value *string, _ z.Ctx) error {
	*value = strings.Join(strings.Fields(*value), " ")
	return nil
}

func IsSupportedCurrency(value *string, _ z.Ctx) bool {
	return common.IsSupportedCurrency(*value)
}
