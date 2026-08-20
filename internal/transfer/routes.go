package transfer

import "github.com/gin-gonic/gin"

func (handler *Handler) RegisterRoutes(router gin.IRoutes) {
	router.POST("/transfers", handler.createTransfer)
}
