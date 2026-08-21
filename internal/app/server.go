// Package app loads application configuration and composes HTTP features.
package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"simplebank/internal/account"
	"simplebank/internal/transfer"
	"simplebank/internal/user"
)

// Store combines the database capabilities required by every HTTP feature.
type Store interface {
	account.Store
	transfer.Store
	user.Store
}

// NewHTTPHandler wires every feature into one router without opening a port.
func NewHTTPHandler(store Store) http.Handler {
	router := gin.Default()
	_ = router.SetTrustedProxies(nil)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	accountService := account.NewService(store)
	account.NewHandler(accountService).RegisterRoutes(router)

	userService := user.NewService(store)
	user.NewHandler(userService).RegisterRoutes(router)

	transferService := transfer.NewService(store, accountService)
	transfer.NewHandler(transferService).RegisterRoutes(router)

	return router
}
