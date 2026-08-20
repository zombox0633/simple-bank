// Package app loads application configuration and composes HTTP features.
package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"simplebank/internal/account"
	"simplebank/internal/httpapi"
	"simplebank/internal/transfer"
)

// Store combines the database capabilities required by every HTTP feature.
type Store interface {
	account.Store
	transfer.Store
}

type Server struct {
	router *gin.Engine
}

func NewServer(store Store) *Server {
	httpapi.RegisterValidators()

	router := gin.Default()
	_ = router.SetTrustedProxies(nil)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	accountService := account.NewService(store)
	account.NewHandler(accountService).RegisterRoutes(router)

	transferService := transfer.NewService(store, accountService)
	transfer.NewHandler(transferService).RegisterRoutes(router)

	return &Server{router: router}
}

// Handler exposes the router for tests and net/http without opening a TCP port.
func (server *Server) Handler() http.Handler {
	return server.router
}
