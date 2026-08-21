package account

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"simplebank/internal/common"
	"simplebank/internal/httpapi"
)

func TestAccountRequestSchemas(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/accounts",
		strings.NewReader(`{"owner":" Alice ","currency":"thb"}`),
	)
	createRequest.Header.Set("Content-Type", "application/json")
	createContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	createContext.Request = createRequest
	var create createAccountRequest
	require.True(t, httpapi.BindJSON(createContext, createAccountSchema, &create))
	require.Equal(t, "alice", create.Owner)
	require.Equal(t, common.CurrencyTHB, create.Currency)

	listRequest := httptest.NewRequest(
		http.MethodGet,
		"/accounts?owner=Alice&page_id=2&page_size=25",
		nil,
	)
	listContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	listContext.Request = listRequest
	var list listAccountsRequest
	require.True(t, httpapi.BindQuery(listContext, listAccountsSchema, &list))
	require.Equal(t, "alice", list.Owner)
	require.Equal(t, int32(2), list.PageID)
	require.Equal(t, int32(25), list.PageSize)

	invalidRequest := httptest.NewRequest(
		http.MethodPost,
		"/accounts",
		strings.NewReader(`{"owner":"alice","currency":"BTC"}`),
	)
	invalidRequest.Header.Set("Content-Type", "application/json")
	invalidContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	invalidContext.Request = invalidRequest
	require.False(t, httpapi.BindJSON(invalidContext, createAccountSchema, &createAccountRequest{}))

	overflowRequest := httptest.NewRequest(
		http.MethodGet,
		"/accounts?owner=alice&page_id=2147483647&page_size=100",
		nil,
	)
	overflowRecorder := httptest.NewRecorder()
	overflowContext, _ := gin.CreateTestContext(overflowRecorder)
	overflowContext.Request = overflowRequest
	require.False(t, httpapi.BindQuery(
		overflowContext,
		listAccountsSchema,
		&listAccountsRequest{},
	))
	require.Equal(t, http.StatusBadRequest, overflowRecorder.Code)
}
