package account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/govalues/decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	db "simplebank/db/sqlc"
)

func TestGetAccountAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	wantAccount := db.Account{
		ID:        42,
		Owner:     "alice",
		Balance:   decimal.MustParse("100.2500"),
		Currency:  "THB",
		CreatedAt: pgtype.Timestamptz{Time: time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC), Valid: true},
	}

	tests := []struct {
		name       string
		accountID  int64
		buildStore func(*testing.T) *stubStore
		wantStatus int
		wantCode   string
	}{
		{
			name:      "OK",
			accountID: wantAccount.ID,
			buildStore: func(_ *testing.T) *stubStore {
				return &stubStore{getAccount: func(context.Context, int64) (db.Account, error) {
					return wantAccount, nil
				}}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:      "NotFound",
			accountID: wantAccount.ID,
			buildStore: func(_ *testing.T) *stubStore {
				return &stubStore{getAccount: func(context.Context, int64) (db.Account, error) {
					return db.Account{}, pgx.ErrNoRows
				}}
			},
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
		{
			name:      "InternalError",
			accountID: wantAccount.ID,
			buildStore: func(_ *testing.T) *stubStore {
				return &stubStore{getAccount: func(context.Context, int64) (db.Account, error) {
					return db.Account{}, errors.New("connection error")
				}}
			},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_ERROR",
		},
		{
			name:       "InvalidID",
			accountID:  0,
			buildStore: func(*testing.T) *stubStore { return &stubStore{} },
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_REQUEST",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := test.buildStore(t)

			router := gin.New()
			NewHandler(NewService(store)).RegisterRoutes(router)

			request := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/accounts/%d", test.accountID),
				nil,
			)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			require.Equal(t, test.wantStatus, response.Code)
			if test.wantStatus == http.StatusOK {
				wantJSON, err := json.Marshal(wantAccount)
				require.NoError(t, err)
				require.JSONEq(t, string(wantJSON), response.Body.String())
				return
			}

			var body struct {
				Code string `json:"code"`
			}
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
			require.Equal(t, test.wantCode, body.Code)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	want := db.Account{ID: 42, Owner: "alice", Currency: "THB"}
	store := &stubStore{
		createAccount: func(_ context.Context, params db.CreateAccountParams) (db.Account, error) {
			require.Equal(t, want.Owner, params.Owner)
			require.Equal(t, want.Currency, params.Currency)
			require.True(t, params.Balance.IsZero())
			return want, nil
		},
	}
	router := gin.New()
	NewHandler(NewService(store)).RegisterRoutes(router)
	request := httptest.NewRequest(
		http.MethodPost,
		"/accounts",
		strings.NewReader(`{"owner":" Alice ","currency":"thb"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusCreated, response.Code)
	wantJSON, err := json.Marshal(want)
	require.NoError(t, err)
	require.JSONEq(t, string(wantJSON), response.Body.String())
}

func TestListAccountsAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	want := []db.Account{{ID: 42, Owner: "alice", Currency: "THB"}}
	store := &stubStore{
		listAccounts: func(_ context.Context, params db.ListAccountsParams) ([]db.Account, error) {
			require.Equal(t, "alice", params.Owner)
			require.Equal(t, int32(25), params.Limit)
			require.Equal(t, int32(50), params.Offset)
			return want, nil
		},
	}
	router := gin.New()
	NewHandler(NewService(store)).RegisterRoutes(router)
	request := httptest.NewRequest(
		http.MethodGet,
		"/accounts?owner=Alice&page_id=3&page_size=25",
		nil,
	)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	wantJSON, err := json.Marshal(want)
	require.NoError(t, err)
	require.JSONEq(t, string(wantJSON), response.Body.String())
}
