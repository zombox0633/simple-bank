package transfer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"simplebank/internal/api"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) createTransfer(ctx *gin.Context) {
	var request createTransferRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		api.WriteRequestError(ctx, err)
		return
	}

	result, err := handler.service.CreateTransfer(
		ctx.Request.Context(),
		request.FromAccountID,
		request.ToAccountID,
		request.Amount,
		request.Currency,
	)
	if err != nil {
		api.WriteError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
