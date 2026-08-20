package transfer

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/govalues/decimal"

	"simplebank/internal/api"
)

type Handler struct {
	service *Service
}

type createTransferResponse struct {
	TransferID    int64           `json:"transfer_id"`
	FromAccountID int64           `json:"from_account_id"`
	ToAccountID   int64           `json:"to_account_id"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	FromBalance   decimal.Decimal `json:"from_balance"`
	CreatedAt     time.Time       `json:"created_at"`
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

	ctx.JSON(http.StatusOK, createTransferResponse{
		TransferID:    result.Transfer.ID,
		FromAccountID: result.Transfer.FromAccountID,
		ToAccountID:   result.Transfer.ToAccountID,
		Amount:        result.Transfer.Amount,
		Currency:      request.Currency,
		FromBalance:   result.FromAccount.Balance,
		CreatedAt:     result.Transfer.CreatedAt.Time,
	})
}
