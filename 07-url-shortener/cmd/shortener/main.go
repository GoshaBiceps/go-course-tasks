package main

import (
	"context"
	"log"
	"net/http"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		"postgres://postgres:postgres@localhost:5432/shortener",
	)

	if err != nil {
		log.Fatal("failed to create datbase pool:", err)
	}

	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("failed to connect to PostgreSQL;", err)
	}

	repo := repository.NewLinkRepository(pool)

	serv := service.NewLinkService(repo)

	hand := handler.NewLinkHandler(serv)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", hand.Shorten)
	mux.HandleFunc("GET /{code}", hand.Redirect)
	mux.HandleFunc("GET / stats/{code}", hand.Stats)

	log.Println("HTTP server starting on :8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("failed to start HTTP server: ", err)
	}
}
