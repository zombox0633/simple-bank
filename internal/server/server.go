// Package server ประกอบ HTTP server และเชื่อมแต่ละ feature เข้ากับ Gin router
package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"simplebank/internal/account"
)

// Store รวมความสามารถด้านฐานข้อมูลที่ทุก feature บน HTTP server ต้องใช้
// ตอนนี้มีเฉพาะ account และสามารถ embed interface ของ feature อื่นเพิ่มภายหลัง
type Store interface {
	account.Store
}

type Server struct {
	router *gin.Engine
}

func New(store Store) *Server {
	router := gin.Default()
	_ = router.SetTrustedProxies(nil)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	accountService := account.NewService(store)
	account.NewHandler(accountService).RegisterRoutes(router)

	return &Server{router: router}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Handler เปิด http.Handler สำหรับทดสอบ routes โดยไม่ต้องเปิด TCP port จริง
func (server *Server) Handler() http.Handler {
	return server.router
}
