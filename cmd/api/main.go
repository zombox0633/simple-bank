package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // register pgx with database/sql

	db "simplebank/db/sqlc"
	"simplebank/internal/config"
	"simplebank/internal/server"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(appConfig.DBDriver, appConfig.DBSource)
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

	if err := apiServer.Start(appConfig.ServerAddress); err != nil {
		log.Fatal("cannot start server 🙀: ", err)
	}
}
