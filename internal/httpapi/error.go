// Package httpapi contains shared Gin validation and HTTP error responses.
package httpapi

import (
	"errors"
	"log"
	"net/http"
	"slices"

	z "github.com/Oudwins/zog"
	"github.com/gin-gonic/gin"

	"simplebank/internal/common"
)

type errorResponse struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func WriteRequestIssues(ctx *gin.Context, issues z.ZogIssueList) {
	slices.SortFunc(issues, func(left *z.ZogIssue, right *z.ZogIssue) int {
		return slices.Compare(left.Path, right.Path)
	})
	for _, issue := range issues {
		switch issue.Code {
		case IssueInvalidUsername:
			WriteError(ctx, common.ErrInvalidUsername)
			return
		case IssueInvalidFullName:
			WriteError(ctx, common.ErrInvalidFullName)
			return
		case IssueInvalidEmail:
			WriteError(ctx, common.ErrInvalidEmail)
			return
		case IssueInvalidPassword:
			WriteError(ctx, common.ErrInvalidPassword)
			return
		case IssueInvalidAmount:
			WriteError(ctx, common.ErrInvalidAmount)
			return
		case IssueUnsupportedCurrency:
			WriteError(ctx, common.ErrUnsupportedCurrency)
			return
		case IssueSameAccount:
			WriteError(ctx, common.ErrSameAccount)
			return
		}
	}

	writeError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
}

// WriteError turns reusable application errors into one consistent HTTP shape.
func WriteError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, common.ErrInvalidUsername):
		writeError(ctx, http.StatusBadRequest, "INVALID_USERNAME", common.ErrInvalidUsername.Error())
	case errors.Is(err, common.ErrInvalidFullName):
		writeError(ctx, http.StatusBadRequest, "INVALID_FULL_NAME", common.ErrInvalidFullName.Error())
	case errors.Is(err, common.ErrInvalidEmail):
		writeError(ctx, http.StatusBadRequest, "INVALID_EMAIL", common.ErrInvalidEmail.Error())
	case errors.Is(err, common.ErrUnsupportedCurrency):
		writeError(ctx, http.StatusBadRequest, "UNSUPPORTED_CURRENCY", "currency must be one of USD, EUR, or THB")
	case errors.Is(err, common.ErrCurrencyMismatch):
		writeError(ctx, http.StatusBadRequest, "CURRENCY_MISMATCH", common.ErrCurrencyMismatch.Error())
	case errors.Is(err, common.ErrInvalidAmount):
		writeError(ctx, http.StatusBadRequest, "INVALID_AMOUNT", common.ErrInvalidAmount.Error())
	case errors.Is(err, common.ErrInvalidPassword):
		writeError(ctx, http.StatusBadRequest, "INVALID_PASSWORD", common.ErrInvalidPassword.Error())
	case errors.Is(err, common.ErrSameAccount):
		writeError(ctx, http.StatusBadRequest, "SAME_ACCOUNT", common.ErrSameAccount.Error())
	case errors.Is(err, common.ErrInsufficientBalance):
		writeError(ctx, http.StatusBadRequest, "INSUFFICIENT_BALANCE", common.ErrInsufficientBalance.Error())
	case errors.Is(err, common.ErrBalanceLimitExceeded):
		writeError(ctx, http.StatusBadRequest, "BALANCE_LIMIT_EXCEEDED", common.ErrBalanceLimitExceeded.Error())
	case errors.Is(err, common.ErrNotFound):
		writeError(ctx, http.StatusNotFound, "NOT_FOUND", common.ErrNotFound.Error())
	case errors.Is(err, common.ErrInvalidReference):
		writeError(ctx, http.StatusBadRequest, "INVALID_REFERENCE", common.ErrInvalidReference.Error())
	case errors.Is(err, common.ErrConflict):
		writeError(ctx, http.StatusConflict, "CONFLICT", common.ErrConflict.Error())
	default:
		log.Printf("request failed: %v", err)
		writeError(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

func writeError(ctx *gin.Context, status int, code string, message string) {
	ctx.JSON(status, errorResponse{Code: code, Error: message})
}
