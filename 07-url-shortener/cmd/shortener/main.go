package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"url-shortener/internal/database"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
	"url-shortener/pkg/cache"
	"url-shortener/pkg/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	redisClient := cache.NewRedisClient(cfg.RedisAddr)

	defer redisClient.Close()

	err = redisClient.Ping(ctx).Err()
	if err != nil {
		slog.Warn(
			"failed to connect to Redis",
			"error", err,
		)
	}

	repo := repository.NewLinkRepository(pool)

	serv := service.NewLinkService(repo, redisClient)

	hand := handler.NewLinkHandler(serv, cfg.BaseURL)

	mux := handler.NewRouter(hand)

	log.Println("HTTP server starting on :8080")

	err = http.ListenAndServe(":"+cfg.HTTPPort, mux)
	if err != nil {
		log.Fatal("failed to start HTTP server: ", err)
	}
}
