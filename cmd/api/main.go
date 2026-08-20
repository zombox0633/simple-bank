package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pgxdecimal "github.com/ColeBurch/pgx-govalues-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	db "simplebank/db/sqlc"
	"simplebank/internal/config"
	"simplebank/internal/server"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	poolConfig, err := pgxpool.ParseConfig(appConfig.DBSource)
	if err != nil {
		log.Fatal("cannot parse database config: ", err)
	}
	poolConfig.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("cannot create database pool: ", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("cannot ping db 😹: ", err)
	}
	log.Println("connected to db 😻")

	store := db.NewStore(pool)
	apiServer := server.New(store)
	httpServer := &http.Server{
		Addr:              appConfig.ServerAddress,
		Handler:           apiServer.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverError := make(chan error, 1)
	log.Printf("server listening on %s", appConfig.ServerAddress)
	go func() {
		serverError <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("cannot start server 🙀: ", err)
		}
	case <-shutdownSignal.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("cannot gracefully stop server: %v", err)
			_ = httpServer.Close()
		}

		if err := <-serverError; err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server stopped with error: %v", err)
		}
	}
}
