package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	db "simplebank/db/sqlc"
)

func TestCreateUserAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Date(2026, time.August, 21, 10, 0, 0, 0, time.UTC)
	localMeanTime := time.FixedZone("LMT", 6*60*60+42*60+4)
	initialPasswordChangedAt := time.Date(1, time.January, 1, 6, 42, 4, 0, localMeanTime)
	wantUser := db.User{
		Username:          "charlie_01",
		FullName:          "Charlie Example",
		Email:             "charlie@example.com",
		PasswordChangedAt: pgtype.Timestamptz{Time: initialPasswordChangedAt, Valid: true},
		CreatedAt:         pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:         pgtype.Timestamptz{Time: now, Valid: true},
	}

	tests := []struct {
		name       string
		body       string
		store      *stubStore
		wantStatus int
		wantCode   string
	}{
		{
			name: "created",
			body: `{
				"username":"Charlie_01",
				"password":"Secret123!",
				"full_name":"  Charlie   Example  ",
				"email":" Charlie@EXAMPLE.COM "
			}`,
			store: &stubStore{
				createUser: func(_ context.Context, arg db.CreateUserParams) (db.User, error) {
					require.Equal(t, wantUser.Username, arg.Username)
					require.Equal(t, wantUser.FullName, arg.FullName)
					require.Equal(t, wantUser.Email, arg.Email)
					createdUser := wantUser
					createdUser.Password = arg.Password
					return createdUser, nil
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid password",
			body: `{
				"username":"charlie_01",
				"password":"short",
				"full_name":"Charlie Example",
				"email":"charlie@example.com"
			}`,
			store:      &stubStore{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_PASSWORD",
		},
		{
			name: "invalid username",
			body: `{
				"username":"_charlie",
				"password":"Secret123!",
				"full_name":"Charlie Example",
				"email":"charlie@example.com"
			}`,
			store:      &stubStore{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_USERNAME",
		},
		{
			name: "invalid full name",
			body: `{
				"username":"charlie_01",
				"password":"Secret123!",
				"full_name":"Charlie123",
				"email":"charlie@example.com"
			}`,
			store:      &stubStore{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_FULL_NAME",
		},
		{
			name: "invalid email",
			body: `{
				"username":"charlie_01",
				"password":"Secret123!",
				"full_name":"Charlie Example",
				"email":"not-an-email"
			}`,
			store:      &stubStore{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_EMAIL",
		},
		{
			name: "duplicate username or email",
			body: `{
				"username":"charlie_01",
				"password":"Secret123!",
				"full_name":"Charlie Example",
				"email":"charlie@example.com"
			}`,
			store: &stubStore{
				createUser: func(_ context.Context, _ db.CreateUserParams) (db.User, error) {
					return db.User{}, &pgconn.PgError{Code: db.SQLStateUniqueViolation}
				},
			},
			wantStatus: http.StatusConflict,
			wantCode:   "CONFLICT",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.New()
			NewHandler(NewService(test.store)).RegisterRoutes(router)

			request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(test.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			require.Equal(t, test.wantStatus, response.Code)
			if test.wantStatus == http.StatusCreated {
				require.NotContains(t, response.Body.String(), `"password":`)
				var body createUserResponse
				require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
				require.Equal(t, wantUser.Username, body.Username)
				require.Equal(t, wantUser.Email, body.Email)
				require.True(t, body.PasswordChangedAt.IsZero())
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
