package account

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"simplebank/internal/httpapi"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) getAccount(ctx *gin.Context) {
	var request getAccountRequest
	if !httpapi.BindURI(ctx, getAccountSchema, &request) {
		return
	}

	account, err := handler.service.GetAccount(ctx.Request.Context(), request.ID)
	if err != nil {
		httpapi.WriteError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (handler *Handler) listAccounts(ctx *gin.Context) {
	var request listAccountsRequest
	if !httpapi.BindQuery(ctx, listAccountsSchema, &request) {
		return
	}

	accounts, err := handler.service.ListAccounts(
		ctx.Request.Context(),
		request.Owner,
		request.PageID,
		request.PageSize,
	)
	if err != nil {
		httpapi.WriteError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (handler *Handler) createAccount(ctx *gin.Context) {
	var request createAccountRequest
	if !httpapi.BindJSON(ctx, createAccountSchema, &request) {
		return
	}

	account, err := handler.service.CreateAccount(
		ctx.Request.Context(),
		request.Owner,
		request.Currency,
	)
	if err != nil {
		httpapi.WriteError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, account)
}
