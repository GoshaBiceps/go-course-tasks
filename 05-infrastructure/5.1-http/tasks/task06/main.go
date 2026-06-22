// Задание: Интеграционная мини-задача — Token Service API
//
// Собери мини-API с тремя эндпоинтами, JSON-хелперами,
// таймаутами и graceful shutdown.
//
// Эндпоинты:
//   GET  /health                — статус 200, {"status":"ok"}
//   POST /api/v1/tokens         — принимает {"user_id":"..."}, возвращает {"token_id":"...","user_id":"..."}
//   GET  /api/v1/tokens/{id}   — возвращает мок {"token_id":"...","user_id":"mock-user","expires_in":3600}
//
// Требования:
//   - единый writeJSON helper
//   - единый writeError helper
//   - маршруты в стиле Go 1.22+
//   - сервер с таймаутами и graceful shutdown

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type errorResponse struct {
	Error string `json:"error"`
}

type healthResponse struct {
	Status string `json:"status"`
}

type createTokenRequest struct {
	UserID string `json:"user_id"`
}

type createTokenResponse struct {
	TokenID string `json:"token_id"`
	UserID  string `json:"user_id"`
}

type getTokenResponse struct {
	TokenID   string `json:"token_id"`
	UserID    string `json:"user_id"`
	ExpiresIn int    `json:"expires_in"`
}

// writeJSON отправляет JSON-ответ с заданным статусом.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// TODO: реализуй writeError(w http.ResponseWriter, status int, msg string)
// должна вызывать writeJSON с телом {"error":"<msg>"}
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	// TODO: верни {"status":"ok"} со статусом 200
	writeJSON(w, http.StatusOK, healthResponse{
		Status: "ok",
	})
}

func handleCreateToken(w http.ResponseWriter, r *http.Request) {
	// TODO: декодируй {"user_id":"..."} из тела запроса
	// При ошибке декодирования — 400, {"error":"invalid json"}
	// Если user_id пустой — 400, {"error":"user_id is required"}
	// Иначе — 201, {"token_id":"tok-<timestamp>","user_id":"<user_id>"}

	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	if req.UserID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "user_id is required"})
		return
	}

	tokenID := fmt.Sprintf("tok-%d", time.Now().UnixNano())

	resp := createTokenResponse{
		TokenID: tokenID,
		UserID:  req.UserID,
	}

	writeJSON(w, http.StatusCreated, resp)
}

func handleGetToken(w http.ResponseWriter, r *http.Request) {
	// TODO: достань id через r.PathValue("id")
	// Верни мок: {"token_id":"<id>","user_id":"mock-user","expires_in":3600}
	id := r.PathValue("id")
	resp := getTokenResponse{
		TokenID:   id,
		UserID:    "mock-user",
		ExpiresIn: 3600,
	}
	writeJSON(w, http.StatusOK, resp)
}

func main() {
	mux := http.NewServeMux()

	// TODO: зарегистрируй обработчики:
	// mux.HandleFunc("GET /health", handleHealth)
	// mux.HandleFunc("POST /api/v1/tokens", handleCreateToken)
	// mux.HandleFunc("GET /api/v1/tokens/{id}", handleGetToken)
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /api/v1/tokens", handleCreateToken)
	mux.HandleFunc("GET /api/v1/tokens/{ud}", handleGetToken)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Println("server started on:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("forced shutdown:", err)
	}
	log.Println("server stopped")
}
