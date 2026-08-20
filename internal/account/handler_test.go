package account

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	db "simplebank/db/sqlc"
)

func TestGetAccountAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	wantAccount := db.Account{
		ID:        42,
		Owner:     "alice",
		Balance:   decimal.RequireFromString("100.2500"),
		Currency:  "THB",
		CreatedAt: time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name       string
		accountID  int64
		buildMock  func(*MockStore)
		wantStatus int
		wantCode   string
	}{
		{
			name:      "OK",
			accountID: wantAccount.ID,
			buildMock: func(store *MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), wantAccount.ID).
					Times(1).
					Return(wantAccount, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:      "NotFound",
			accountID: wantAccount.ID,
			buildMock: func(store *MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), wantAccount.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
		{
			name:      "InternalError",
			accountID: wantAccount.ID,
			buildMock: func(store *MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), wantAccount.ID).
					Times(1).
					Return(db.Account{}, errors.New("connection error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_ERROR",
		},
		{
			name:       "InvalidID",
			accountID:  0,
			buildMock:  func(*MockStore) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_REQUEST",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := NewMockStore(gomock.NewController(t))
			test.buildMock(store)

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
