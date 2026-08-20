package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"

	"simplebank/internal/common"
)

type validationRequest struct {
	FromAccountID int64           `binding:"required,min=1"`
	ToAccountID   int64           `binding:"required,min=1,nefield=FromAccountID"`
	Amount        decimal.Decimal `binding:"money"`
	Currency      string          `binding:"required,currency"`
}

func TestCustomValidatorsAndRequestErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	RegisterValidators()

	validRequest := validationRequest{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        decimal.MustNew(10, 0),
		Currency:      common.CurrencyTHB,
	}
	require.NoError(t, binding.Validator.ValidateStruct(validRequest))

	tests := []struct {
		name     string
		request  validationRequest
		wantCode string
	}{
		{
			name: "invalid amount",
			request: validationRequest{
				FromAccountID: 1, ToAccountID: 2,
				Amount: decimal.MustNew(-10, 0), Currency: common.CurrencyTHB,
			},
			wantCode: "INVALID_AMOUNT",
		},
		{
			name: "unsupported currency",
			request: validationRequest{
				FromAccountID: 1, ToAccountID: 2,
				Amount: decimal.MustNew(10, 0), Currency: "BTC",
			},
			wantCode: "UNSUPPORTED_CURRENCY",
		},
		{
			name: "same account",
			request: validationRequest{
				FromAccountID: 1, ToAccountID: 1,
				Amount: decimal.MustNew(10, 0), Currency: common.CurrencyTHB,
			},
			wantCode: "SAME_ACCOUNT",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := binding.Validator.ValidateStruct(test.request)
			require.Error(t, err)

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			WriteRequestError(ctx, err)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			var response errorResponse
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
			require.Equal(t, test.wantCode, response.Code)
		})
	}
}
