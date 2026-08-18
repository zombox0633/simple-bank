package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver แบบ database/sql → ทำให้ sql.Open("pgx", ...) ใช้ได้
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	conn, err := sql.Open("pgx", os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("cannot ping db 😹 :", err)
	}
	log.Println("connected to db 😻")

	// ต่อไป (phase 4): store := db.New(conn) แล้วเอา store ไปใส่ใน handler

	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = ":8080" // fallback เผื่อ .env ไม่มีค่านี้
	}

	r := gin.Default()
	_ = r.SetTrustedProxies(nil) // แอปไม่ได้อยู่หลัง proxy → ปิด trusted proxies

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	//test ci

	if err := r.Run(addr); err != nil {
		log.Fatal("cannot start server 🙀 : ", err)
	}
}
