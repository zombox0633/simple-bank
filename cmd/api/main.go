package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // register pgx with database/sql
	"github.com/joho/godotenv"

	db "simplebank/db/sqlc"
	"simplebank/internal/server"
)

func main() {
	_ = godotenv.Load()

	conn, err := sql.Open("pgx", os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("cannot ping db 😹: ", err)
	}
	log.Println("connected to db 😻")

	store := db.NewStore(conn)
	apiServer := server.New(store)

	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = ":8080"
	}

	if err := apiServer.Start(address); err != nil {
		log.Fatal("cannot start server 🙀: ", err)
	}
}
