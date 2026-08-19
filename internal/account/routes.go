package account

import "github.com/gin-gonic/gin"

func (handler *Handler) RegisterRoutes(router gin.IRoutes) {
	router.GET("/accounts/:id", handler.getAccount)
	router.GET("/accounts", handler.listAccounts)
	router.POST("/accounts", handler.createAccount)
}
