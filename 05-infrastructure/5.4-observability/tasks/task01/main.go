package main

import (
	"errors"
	"log/slog"
	"os"
)

// TODO: Implement newLogger(service, env string) *slog.Logger
// - Use slog.NewJSONHandler writing to os.Stdout
// - Set level to slog.LevelInfo
// - Attach base fields: "service" and "env"

func newLogger(service, env string) *slog.Logger {

	handler := slog.NewJSONHandler( //обработчик логов
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	logger := slog.New(handler).With( // создали логер он же пишет событие
		"service", service,
		"env", env,
	)

	return logger
}

func main() {
	logger := newLogger("token-service", "development")

	// TODO: Log service startup at Info level with a "version" field
	logger.Info(
		"service started",
		"version",
		"v1.0.0",
	)

	// TODO: Log a simulated incoming request at Info level
	// with fields: method, path, request_id

	logger.Info(
		"incoming request",
		"method", "GET",
		"path", "/api/v1/token",
		"request_id", "req-123",
	)

	// TODO: Log a simulated error at Error level
	// using errors.New("connection timeout") as the error field
	err := errors.New("connection timeout")

	logger.Error(
		"database request failed",
		"error", err,
	)

}
