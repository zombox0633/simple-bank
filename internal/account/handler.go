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
	if err := ctx.ShouldBindUri(&request); err != nil {
		httpapi.WriteRequestError(ctx, err)
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
	if err := ctx.ShouldBindQuery(&request); err != nil {
		httpapi.WriteRequestError(ctx, err)
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
	if err := ctx.ShouldBindJSON(&request); err != nil {
		httpapi.WriteRequestError(ctx, err)
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
