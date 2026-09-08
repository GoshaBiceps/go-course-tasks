// Задание: Интеграционная мини-задача — HTTP API → gRPC TokenService
//
// Реализуй HTTP-сервис аутентификации, который внутри вызывает gRPC TokenService.
//
// HTTP API (Auth Service):
//   POST /auth/login   {"user_id":"u-1"} → {"token":"tok-1","user_id":"u-1"}
//   POST /auth/logout  {"token":"tok-1"} → {"revoked":true}
//   GET  /auth/verify  Header: Authorization: Bearer tok-1 → {"user_id":"u-1","valid":true}
//
// Внутри HTTP-хендлеры вызывают TokenServiceClient (in-process или через grpc.Dial).
//
// Требования:
//   - Передавай deadline из HTTP-запроса в gRPC-вызов (ctx с таймаутом 2s)
//   - gRPC-статус PermissionDenied → HTTP 403
//   - gRPC-статус NotFound        → HTTP 404
//   - gRPC-статус InvalidArgument → HTTP 400
//   - Логируй каждый запрос (метод, путь, статус)
//
// Ожидаемый результат:
//   $ curl -X POST http://localhost:8080/auth/login -d '{"user_id":"u-1"}'
//   {"token":"tok-1","user_id":"u-1"}
//
//   $ curl http://localhost:8080/auth/verify -H "Authorization: Bearer tok-1"
//   {"user_id":"u-1","valid":true}
//
//   $ curl -X POST http://localhost:8080/auth/logout -d '{"token":"tok-1"}'
//   {"revoked":true}
//
//   $ curl http://localhost:8080/auth/verify -H "Authorization: Bearer tok-1"
//   {"error":"permission denied"}  HTTP 403

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// --- gRPC-статусы (без внешних зависимостей) ---

type StatusCode int

const (
	CodeInvalidArgument  StatusCode = 3
	CodeNotFound         StatusCode = 5
	CodePermissionDenied StatusCode = 7
)

type StatusError struct {
	Code    StatusCode
	Message string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("rpc error: code = %d desc = %s", e.Code, e.Message)
}

// --- TokenService (in-process, заменяет реальный gRPC-сервер) ---

type IssueTokenRequest struct{ UserID string }
type IssueTokenResponse struct{ TokenID, UserID string }
type ValidateTokenRequest struct{ TokenID string }
type ValidateTokenResponse struct {
	UserID string
	Valid  bool
}
type RevokeTokenRequest struct{ TokenID string }
type RevokeTokenResponse struct{ Revoked bool }

type TokenServiceClient interface {
	IssueToken(ctx context.Context, req *IssueTokenRequest) (*IssueTokenResponse, error)
	ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error)
	RevokeToken(ctx context.Context, req *RevokeTokenRequest) (*RevokeTokenResponse, error)
}

type inProcessTokenService struct {
	tokens map[string]*tokenRecord
	nextID int
}

type tokenRecord struct {
	userID  string
	revoked bool
}

func newTokenService() TokenServiceClient {
	return &inProcessTokenService{tokens: make(map[string]*tokenRecord), nextID: 1}
}

func (s *inProcessTokenService) IssueToken(_ context.Context, req *IssueTokenRequest) (*IssueTokenResponse, error) {
	if req.UserID == "" {
		return nil, &StatusError{Code: CodeInvalidArgument, Message: "user_id is required"}
	}
	id := fmt.Sprintf("tok-%d", s.nextID)
	s.nextID++
	s.tokens[id] = &tokenRecord{userID: req.UserID}
	return &IssueTokenResponse{TokenID: id, UserID: req.UserID}, nil
}

func (s *inProcessTokenService) ValidateToken(_ context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	t, ok := s.tokens[req.TokenID]
	if !ok {
		return nil, &StatusError{Code: CodeNotFound, Message: "token not found"}
	}
	if t.revoked {
		return nil, &StatusError{Code: CodePermissionDenied, Message: "token revoked"}
	}
	return &ValidateTokenResponse{UserID: t.userID, Valid: true}, nil
}

func (s *inProcessTokenService) RevokeToken(_ context.Context, req *RevokeTokenRequest) (*RevokeTokenResponse, error) {
	t, ok := s.tokens[req.TokenID]
	if !ok {
		return nil, &StatusError{Code: CodeNotFound, Message: "token not found"}
	}
	t.revoked = true
	return &RevokeTokenResponse{Revoked: true}, nil
}

// --- HTTP helpers ---

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// grpcStatusToHTTP конвертирует gRPC-статус в HTTP-код.
// TODO: реализуй маппинг:
//
//	CodeInvalidArgument  → 400
//	CodeNotFound         → 404
//	CodePermissionDenied → 403
//	иначе                → 500

func grpcStatusToHTTP(code StatusCode) int {
	switch code {

	case CodeInvalidArgument:
		return http.StatusBadRequest

	case CodeNotFound:
		return http.StatusNotFound

	case CodePermissionDenied:
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}

// writeGRPCError отправляет JSON-ошибку с правильным HTTP-статусом.
func writeGRPCError(w http.ResponseWriter, err error) {
	if se, ok := err.(*StatusError); ok {
		writeJSON(w, grpcStatusToHTTP(se.Code), map[string]string{"error": se.Message})
		return
	}
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
}

// --- HTTP handlers ---

// TODO: реализуй handleLogin(svc TokenServiceClient, logger *slog.Logger) http.HandlerFunc
// Шаги:
//   1. Декодируй JSON: {"user_id":"..."}
//   2. Создай ctx с таймаутом 2s: ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
//   3. Вызови svc.IssueToken(ctx, ...)
//   4. При ошибке → writeGRPCError; при успехе → writeJSON 200 {"token":"...","user_id":"..."}
//   5. Залогируй: logger.Info("login", "user_id", ..., "token_id", ...)

func handleLogin(
	svc TokenServiceClient,
	logger *slog.Logger,
) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		var body struct {
			UserID string `json:"user_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]string{
					"error": "invalid json",
				},
			)
			return
		}

		ctx, cancel := context.WithTimeout(
			r.Context(),
			2*time.Second,
		)
		defer cancel()

		resp, err := svc.IssueToken(
			ctx,
			&IssueTokenRequest{
				UserID: body.UserID,
			},
		)

		if err != nil {
			writeGRPCError(w, err)
			return
		}

		logger.Info(
			"login",
			"user_id", resp.UserID,
			"token_id", resp.TokenID,
		)

		writeJSON(
			w,
			http.StatusOK,
			map[string]string{
				"token":   resp.TokenID,
				"user_id": resp.UserID,
			},
		)
	}
}

// TODO: реализуй handleVerify(svc TokenServiceClient, logger *slog.Logger) http.HandlerFunc
// Шаги:
//   1. Извлеки Bearer-токен из заголовка Authorization
//   2. Если нет → 401 {"error":"missing token"}
//   3. ctx с таймаутом 2s
//   4. Вызови svc.ValidateToken(ctx, ...)
//   5. При ошибке → writeGRPCError; при успехе → writeJSON 200 {"user_id":"...","valid":true}

func handleVerify(
	svc TokenServiceClient,
	logger *slog.Logger,
) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		token, ok := extractBearer(r)

		if !ok {
			writeJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{
					"error": "missing token",
				},
			)
			return
		}

		ctx, cancel := context.WithTimeout(
			r.Context(),
			2*time.Second,
		)
		defer cancel()

		resp, err := svc.ValidateToken(
			ctx,
			&ValidateTokenRequest{
				TokenID: token,
			},
		)

		if err != nil {
			writeGRPCError(w, err)
			return
		}

		logger.Info(
			"verify",
			"user_id", resp.UserID,
			"token_id", token,
		)

		writeJSON(
			w,
			http.StatusOK,
			map[string]any{
				"user_id": resp.UserID,
				"valid":   resp.Valid,
			},
		)
	}
}

// TODO: реализуй handleLogout(svc TokenServiceClient, logger *slog.Logger) http.HandlerFunc
// Шаги:
//   1. Декодируй JSON: {"token":"..."}
//   2. ctx с таймаутом 2s
//   3. Вызови svc.RevokeToken(ctx, ...)
//   4. При ошибке → writeGRPCError; при успехе → writeJSON 200 {"revoked":true}

func handleLogout(
	svc TokenServiceClient, // клиент TokenService
	logger *slog.Logger, // логгер
) http.HandlerFunc {

	// Возвращаем HTTP-хендлер для POST /auth/logout.
	return func(
		w http.ResponseWriter, // сюда пишем HTTP-ответ
		r *http.Request, // HTTP-запрос клиента
	) {

		// Сюда декодируем JSON:
		// {"token":"tok-1"}
		var body struct {
			Token string `json:"token"`
		}

		// Читаем тело запроса и декодируем JSON в body.
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {

			// Если JSON некорректный — HTTP 400.
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]string{
					"error": "invalid json",
				},
			)
			return
		}

		// Берём context текущего HTTP-запроса
		// и добавляем timeout 2 секунды.
		ctx, cancel := context.WithTimeout(
			r.Context(),
			2*time.Second,
		)
		defer cancel()

		// Просим TokenService отозвать токен.
		resp, err := svc.RevokeToken(
			ctx,
			&RevokeTokenRequest{
				TokenID: body.Token,
			},
		)

		// Если TokenService вернул ошибку —
		// переводим gRPC-ошибку в HTTP-ответ.
		if err != nil {
			writeGRPCError(w, err)
			return
		}

		// Логируем успешный logout.
		logger.Info(
			"logout",
			"token_id", body.Token,
		)

		// Возвращаем клиенту:
		// {"revoked":true}
		writeJSON(
			w,
			http.StatusOK,
			map[string]bool{
				"revoked": resp.Revoked,
			},
		)
	}
}

// extractBearer извлекает токен из "Authorization: Bearer <token>".
func extractBearer(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return "", false
	}
	return strings.TrimPrefix(h, "Bearer "), true
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"service", "auth-service",
	)

	svc := newTokenService()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/login", handleLogin(svc, logger))
	mux.HandleFunc("GET /auth/verify", handleVerify(svc, logger))
	mux.HandleFunc("POST /auth/logout", handleLogout(svc, logger))

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	logger.Info("server started", "addr", ":8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", "error", err)
	}
}
