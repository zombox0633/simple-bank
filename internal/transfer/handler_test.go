package transfer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	db "simplebank/db/sqlc"
	"simplebank/internal/api"
	"simplebank/internal/common"
)

type stubStore struct {
	getAccount func(context.Context, int64) (db.Account, error)
	transferTx func(context.Context, db.TransferTxParams) (*db.TransferTxResult, error)
}

func (store *stubStore) GetAccount(ctx context.Context, id int64) (db.Account, error) {
	return store.getAccount(ctx, id)
}

func (store *stubStore) TransferTx(
	ctx context.Context,
	params db.TransferTxParams,
) (*db.TransferTxResult, error) {
	return store.transferTx(ctx, params)
}

func TestCreateTransferAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api.RegisterValidators()

	tests := []struct {
		name       string
		body       string
		store      *stubStore
		wantStatus int
		wantCode   string
		wantTags   []string
	}{
		{
			name: "OK",
			body: `{"from_account_id":1,"to_account_id":2,"amount":"10.1234","currency":"THB"}`,
			store: &stubStore{
				getAccount: func(_ context.Context, id int64) (db.Account, error) {
					return db.Account{
						ID:       id,
						Balance:  decimal.RequireFromString("100.0000"),
						Currency: common.CurrencyTHB,
					}, nil
				},
				transferTx: func(_ context.Context, params db.TransferTxParams) (*db.TransferTxResult, error) {
					require.Equal(t, int64(1), params.FromAccountID)
					require.Equal(t, int64(2), params.ToAccountID)
					require.True(t, params.Amount.Equal(decimal.RequireFromString("10.1234")))
					return &db.TransferTxResult{Transfer: db.Transfer{ID: 99}}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "InvalidParams",
			body:       `{"from_account_id":1,"to_account_id":2,"amount":"1.23456","currency":"BTC"}`,
			store:      &stubStore{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_REQUEST",
			wantTags:   []string{"'money'", "'currency'"},
		},
		{
			name: "CurrencyMismatch",
			body: `{"from_account_id":1,"to_account_id":2,"amount":"10.0000","currency":"THB"}`,
			store: &stubStore{
				getAccount: func(_ context.Context, id int64) (db.Account, error) {
					return db.Account{ID: id, Balance: decimal.NewFromInt(100), Currency: common.CurrencyUSD}, nil
				},
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "CURRENCY_MISMATCH",
		},
		{
			name: "InsufficientBalance",
			body: `{"from_account_id":1,"to_account_id":2,"amount":"10.0000","currency":"THB"}`,
			store: &stubStore{
				getAccount: func(_ context.Context, id int64) (db.Account, error) {
					balance := decimal.NewFromInt(100)
					if id == 1 {
						balance = decimal.NewFromInt(5)
					}
					return db.Account{ID: id, Balance: balance, Currency: common.CurrencyTHB}, nil
				},
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INSUFFICIENT_BALANCE",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.New()
			NewHandler(NewService(test.store, test.store)).RegisterRoutes(router)

			request := httptest.NewRequest(http.MethodPost, "/transfers", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			require.Equal(t, test.wantStatus, response.Code)
			if test.wantStatus == http.StatusOK {
				var result db.TransferTxResult
				require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
				require.Equal(t, int64(99), result.Transfer.ID)
				return
			}

			var errorBody struct {
				Code  string `json:"code"`
				Error string `json:"error"`
			}
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &errorBody))
			require.Equal(t, test.wantCode, errorBody.Code)
			for _, tag := range test.wantTags {
				require.Contains(t, errorBody.Error, tag)
			}
		})
	}
}
