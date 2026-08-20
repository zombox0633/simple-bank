// Package server ประกอบ HTTP server และเชื่อมแต่ละ feature เข้ากับ Gin router
package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"simplebank/internal/account"
	"simplebank/internal/api"
	"simplebank/internal/transfer"
)

// Store รวมความสามารถด้านฐานข้อมูลที่ทุก feature บน HTTP server ต้องใช้
type Store interface {
	account.Store
	transfer.Store
}

type Server struct {
	router *gin.Engine
}

func New(store Store) *Server {
	api.RegisterValidators()

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

// Handler เปิด http.Handler สำหรับทดสอบ routes โดยไม่ต้องเปิด TCP port จริง
func (server *Server) Handler() http.Handler {
	return server.router
}
