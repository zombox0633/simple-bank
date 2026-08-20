package transfer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/govalues/decimal"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	db "simplebank/db/sqlc"
	"simplebank/internal/common"
	"simplebank/internal/httpapi"
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
	httpapi.RegisterValidators()
	createdAt := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		body       string
		store      *stubStore
		wantStatus int
		wantCode   string
	}{
		{
			name: "OK",
			body: `{"from_account_id":1,"to_account_id":2,"amount":"10","currency":"THB"}`,
			store: &stubStore{
				getAccount: func(_ context.Context, id int64) (db.Account, error) {
					return db.Account{
						ID:       id,
						Balance:  decimal.MustParse("100.0000"),
						Currency: common.CurrencyTHB,
					}, nil
				},
				transferTx: func(_ context.Context, params db.TransferTxParams) (*db.TransferTxResult, error) {
					require.Equal(t, int64(1), params.FromAccountID)
					require.Equal(t, int64(2), params.ToAccountID)
					require.True(t, params.Amount.Equal(decimal.MustNew(10, 0)))
					return &db.TransferTxResult{
						Transfer: db.Transfer{
							ID:            99,
							FromAccountID: params.FromAccountID,
							ToAccountID:   params.ToAccountID,
							Amount:        params.Amount,
							CreatedAt:     pgtype.Timestamptz{Time: createdAt, Valid: true},
						},
						FromAccount: db.Account{Balance: decimal.MustNew(90, 0)},
						ToAccount:   db.Account{Balance: decimal.MustNew(110, 0)},
					}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NegativeAmount",
			body:       `{"from_account_id":1,"to_account_id":2,"amount":"-10","currency":"THB"}`,
			store:      &stubStore{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_AMOUNT",
		},
		{
			name:       "UnsupportedCurrency",
			body:       `{"from_account_id":1,"to_account_id":2,"amount":"10","currency":"BTC"}`,
			store:      &stubStore{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "UNSUPPORTED_CURRENCY",
		},
		{
			name: "CurrencyMismatch",
			body: `{"from_account_id":1,"to_account_id":2,"amount":"10.0000","currency":"THB"}`,
			store: &stubStore{
				getAccount: func(_ context.Context, id int64) (db.Account, error) {
					return db.Account{ID: id, Balance: decimal.MustNew(100, 0), Currency: common.CurrencyUSD}, nil
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
					balance := decimal.MustNew(100, 0)
					if id == 1 {
						balance = decimal.MustNew(5, 0)
					}
					return db.Account{ID: id, Balance: balance, Currency: common.CurrencyTHB}, nil
				},
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INSUFFICIENT_BALANCE",
		},
		{
			name: "BalanceLimitExceeded",
			body: `{"from_account_id":1,"to_account_id":2,"amount":"10.0000","currency":"THB"}`,
			store: &stubStore{
				getAccount: func(_ context.Context, id int64) (db.Account, error) {
					balance := decimal.MustNew(100, 0)
					if id == 2 {
						balance = common.MaxMoneyAmount
					}
					return db.Account{ID: id, Balance: balance, Currency: common.CurrencyTHB}, nil
				},
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "BALANCE_LIMIT_EXCEEDED",
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
				var result createTransferResponse
				require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
				require.Equal(t, int64(99), result.TransferID)
				require.Equal(t, int64(1), result.FromAccountID)
				require.Equal(t, int64(2), result.ToAccountID)
				require.True(t, result.Amount.Equal(decimal.MustNew(10, 0)))
				require.Equal(t, common.CurrencyTHB, result.Currency)
				require.True(t, result.FromBalance.Equal(decimal.MustNew(90, 0)))
				require.Equal(t, createdAt, result.CreatedAt)

				var rawResponse map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(response.Body.Bytes(), &rawResponse))
				require.NotContains(t, rawResponse, "transfer")
				require.NotContains(t, rawResponse, "from_account")
				require.NotContains(t, rawResponse, "to_account")
				require.NotContains(t, rawResponse, "from_entry")
				require.NotContains(t, rawResponse, "to_entry")
				return
			}

			var errorBody struct {
				Code  string `json:"code"`
				Error string `json:"error"`
			}
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &errorBody))
			require.Equal(t, test.wantCode, errorBody.Code)
		})
	}
}
