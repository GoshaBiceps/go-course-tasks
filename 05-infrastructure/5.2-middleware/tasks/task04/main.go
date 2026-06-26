// Задание: Auth middleware
//
// Реализуй middleware аутентификации на основе интерфейса TokenVerifier.
// Middleware извлекает Bearer-токен, верифицирует его и кладёт claims в context.
// При отказе — возвращает 401.
//
// Ожидаемый результат:
//   $ curl http://localhost:8080/api/v1/me
//   {"error":"authorization header is empty"}
//
//   $ curl -H "Authorization: Bearer valid-token" http://localhost:8080/api/v1/me
//   {"user_id":"user-123","role":"admin"}
//
//   $ curl -H "Authorization: Bearer bad-token" http://localhost:8080/api/v1/me
//   {"error":"invalid token"}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

// mockVerifier проверяет токен без реальной криптографии (для учебных целей).
type mockVerifier struct{}

func (v *mockVerifier) Verify(token string) (Claims, error) {
	if token == "valid-token" {
		return Claims{UserID: "user-123", Role: "admin"}, nil
	}
	return Claims{}, errors.New("invalid token")
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("authorization header is empty")
	}
	token, ok := strings.CutPrefix(header, "Bearer ") // проверяет начинается ли строка с Bearer и отрезает  его и возвращает остаток
	if !ok || token == "" {
		return "", errors.New("invalid authorization header format")
	}
	return token, nil
}

// TODO: реализуй AuthMiddleware(verifier TokenVerifier) func(http.Handler) http.Handler
// Логика:
//   1. Достать заголовок Authorization из r.Header
//   2. Вызвать extractBearerToken; при ошибке — 401 с {"error":"<msg>"}
//   3. Вызвать verifier.Verify(token); при ошибке — 401 с {"error":"invalid token"}
//   4. Положить claims в context: context.WithValue(r.Context(), claimsKey, claims)
//   5. Передать управление следующему обработчику с обновлённым request

func AuthMiddleware(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := r.Header.Get("Authorization") // достоем занчение из hedera

			token, err := extractBearerToken(raw) // достаем чистый токен
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{
					"error": err.Error(),
				})
				return
			}

			claims, err := verifier.Verify(token) // проверяем токен
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{
					"error": err.Error(),
				})
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: достань claims из context:
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
	verifier := &mockVerifier{}

	mux := http.NewServeMux()
	mux.Handle("GET /api/v1/me", AuthMiddleware(verifier)(http.HandlerFunc(meHandler)))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("server error:", err)
	}
}
