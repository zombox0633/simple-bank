package user

import "github.com/gin-gonic/gin"

func (handler *Handler) RegisterRoutes(router gin.IRoutes) {
	router.POST("/users", handler.createUser)
}
