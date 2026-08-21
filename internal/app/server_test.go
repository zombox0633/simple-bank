package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	db "simplebank/db/sqlc"
)

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	apiHandler := NewHTTPHandler(&db.Store{})
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response := httptest.NewRecorder()

	apiHandler.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.JSONEq(t, `{"message":"pong"}`, response.Body.String())
}
