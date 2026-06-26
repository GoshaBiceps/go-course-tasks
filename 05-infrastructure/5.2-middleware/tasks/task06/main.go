// Задание: Интеграционная мини-задача — публичный и защищённый endpoint
//
// Собери API с двумя маршрутами:
//   GET /health       — публичный, без авторизации
//   GET /api/v1/me    — только с валидным Bearer-токеном
//
// Требования:
//   - цепочка middleware только на защищённом endpoint
//   - структурированное логирование запросов (slog)
//   - единый формат 401 ошибки
//   - замена верификатора без изменения handler'ов
//
// Ожидаемый результат:
//   $ curl http://localhost:8080/health
//   {"status":"ok"}
//
//   $ curl http://localhost:8080/api/v1/me
//   {"error":"authorization header is empty"}
//
//   $ curl -H "Authorization: Bearer valid-token" http://localhost:8080/api/v1/me
//   {"user_id":"user-123","role":"admin"}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Claims struct {
	UserID string
	Role   string
}

type contextKey string

const claimsKey contextKey = "claims"

type TokenVerifier interface {
	Verify(token string) (Claims, error)
}

type mockVerifier struct{}

func (v *mockVerifier) Verify(token string) (Claims, error) {
	if token == "valid-token" {
		return Claims{UserID: "user-123", Role: "admin"}, nil
	}
	return Claims{}, errors.New("invalid token")
}

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type statusRecorder struct { // структура для перехвата статуса
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status                    // статус для нас
	r.ResponseWriter.WriteHeader(status) // предали статус дальше в оригинальный ResponceWriter
}

// TODO: реализуй LoggingMiddleware(logger *slog.Logger) Middleware
// Логируй: method, path, status, duration_ms
func LoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)

			logger.Info(
				"request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", duration.Milliseconds(),
			)
		})
	}
}

// TODO: реализуй AuthMiddleware(verifier TokenVerifier) Middleware
// Извлекай Bearer-токен, верифицируй, клади claims в context
// На ошибке — 401 {"error":"<msg>"}

func extractBearerToken(header string) (string, error) {
	const prefix = "Bearer "

	if !strings.HasPrefix(header, prefix) {
		return "", errors.New("missing bearer token")
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))

	if token == "" {
		return "", errors.New("empty token")
	}

	return token, nil
}

func AuthMiddleware(verifier TokenVerifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := r.Header.Get("Authorization")

			token, err := extractBearerToken(raw) // достали чистый токен
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{
					"error": err.Error(),
				})
				return
			}

			claims, err := verifier.Verify(token)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{
					"error": err.Error(),
				})
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: достань claims из r.Context() и верни {"user_id":"...","role":"..."}
	claims, ok := r.Context().Value(claimsKey).(Claims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "claims not found",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"user_id": claims.UserID,
		"role":    claims.Role,
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	verifier := &mockVerifier{}

	mux := http.NewServeMux()

	// публичный
	mux.Handle(
		"GET /health",
		LoggingMiddleware(logger)(http.HandlerFunc(healthHandler)),
	)

	// защищенный
	protected := Chain(
		http.HandlerFunc(meHandler),
		LoggingMiddleware(logger),
		AuthMiddleware(verifier),
	)

	mux.Handle("GET /api/v1/me", protected)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info("server started", "addr", ":8080")
	_ = verifier // убери после реализации
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server failed", "error", err)
	}
}
