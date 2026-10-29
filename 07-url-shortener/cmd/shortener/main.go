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

	linkRepo := repository.NewLinkRepository(pool)
	//userRepo := repository.NewUserRepository(pool)

	linkService := service.NewLinkService(linkRepo, redisClient)
	//userService := service.NewUserService(userRepo)

	LinkHandler := handler.NewLinkHandler(linkService, cfg.BaseURL)

	mux := handler.NewRouter(LinkHandler)

	log.Println("HTTP server starting on :8080")

	err = http.ListenAndServe(":"+cfg.HTTPPort, mux)
	if err != nil {
		log.Fatal("failed to start HTTP server: ", err)
	}
}
