package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string // адрес подключения
	HTTPPort    string // порт на котором рабоатет сервер
	BaseURL     string // адрес из которого формируются короткие ссылки
}

func Load() (Config, error) {

	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	httpPORT := os.Getenv("HTTP_PORT")

	if httpPORT == "" {
		httpPORT = "8080"
	}

	baseURL := os.Getenv("BASE_URL")

	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	cfg := Config{
		DatabaseURL: dbURL,
		HTTPPort:    httpPORT,
		BaseURL:     baseURL,
	}

	return cfg, nil
}
