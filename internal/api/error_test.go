package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"simplebank/internal/common"
)

func TestWriteError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		err        error
		wantStatus int
		wantCode   string
	}{
		{common.ErrUnsupportedCurrency, http.StatusBadRequest, "UNSUPPORTED_CURRENCY"},
		{common.ErrCurrencyMismatch, http.StatusBadRequest, "CURRENCY_MISMATCH"},
		{common.ErrInvalidAmount, http.StatusBadRequest, "INVALID_AMOUNT"},
		{common.ErrSameAccount, http.StatusBadRequest, "SAME_ACCOUNT"},
		{common.ErrInsufficientBalance, http.StatusBadRequest, "INSUFFICIENT_BALANCE"},
		{common.ErrBalanceLimitExceeded, http.StatusBadRequest, "BALANCE_LIMIT_EXCEEDED"},
		{common.ErrNotFound, http.StatusNotFound, "NOT_FOUND"},
		{common.ErrInvalidReference, http.StatusBadRequest, "INVALID_REFERENCE"},
		{common.ErrConflict, http.StatusConflict, "CONFLICT"},
		{errors.New("private database detail"), http.StatusInternalServerError, "INTERNAL_ERROR"},
	}

	for _, test := range tests {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		WriteError(ctx, test.err)

		require.Equal(t, test.wantStatus, recorder.Code)
		var response errorResponse
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
		require.Equal(t, test.wantCode, response.Code)
	}
}

func TestWriteRequestErrorHidesParserDetails(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	WriteRequestError(ctx, errors.New("private parser detail"))

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	var response errorResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.Equal(t, "INVALID_REQUEST", response.Code)
	require.Equal(t, "invalid request", response.Error)
}
