package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"simplebank/internal/common"
)

type errorResponse struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func WriteRequestError(ctx *gin.Context, err error) {
	writeError(ctx, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
}

// WriteError turns reusable application errors into one consistent HTTP shape.
func WriteError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, common.ErrUnsupportedCurrency):
		writeError(ctx, http.StatusBadRequest, "UNSUPPORTED_CURRENCY", "currency must be one of USD, EUR, or THB")
	case errors.Is(err, common.ErrCurrencyMismatch):
		writeError(ctx, http.StatusBadRequest, "CURRENCY_MISMATCH", common.ErrCurrencyMismatch.Error())
	case errors.Is(err, common.ErrInvalidAmount):
		writeError(ctx, http.StatusBadRequest, "INVALID_AMOUNT", common.ErrInvalidAmount.Error())
	case errors.Is(err, common.ErrSameAccount):
		writeError(ctx, http.StatusBadRequest, "SAME_ACCOUNT", common.ErrSameAccount.Error())
	case errors.Is(err, common.ErrInsufficientBalance):
		writeError(ctx, http.StatusBadRequest, "INSUFFICIENT_BALANCE", common.ErrInsufficientBalance.Error())
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
