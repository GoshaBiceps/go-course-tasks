// Задание: Unary-интерцепторы для gRPC
//
// Реализуй три интерцептора для TokenService:
//   1. LoggingInterceptor   — логирует метод, длительность и ошибку
//   2. AuthInterceptor      — проверяет токен сервис-к-сервису в metadata
//   3. MetricsInterceptor   — счётчик вызовов и ошибок
//
// Интерцепторы должны быть цепочкой (ChainInterceptors).
//
// Ожидаемый результат:
//   [metrics] /token.TokenService/IssueToken calls=1
//   [log] /token.TokenService/IssueToken 1.2ms OK
//   Issued: tok-1
//   [metrics] /token.TokenService/IssueToken calls=2
//   [log] /token.TokenService/IssueToken 0.8ms OK
//   Validated: u-1
//   [metrics] /token.TokenService/ValidateToken calls=1
//   [log] /token.TokenService/ValidateToken 0.5ms PermissionDenied

package main

import (
	"context"
	"fmt"
	"time"
)

// --- Типы интерцепторов ---

type UnaryHandler func(ctx context.Context, req any) (any, error)

type UnaryInterceptor func(ctx context.Context, method string, req any, handler UnaryHandler) (any, error)

// --- gRPC-статусы ---

type StatusCode int

const (
	CodeOK               StatusCode = 0
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

func codeName(code StatusCode) string {
	switch code {
	case CodeOK:
		return "OK"
	case CodeInvalidArgument:
		return "InvalidArgument"
	case CodeNotFound:
		return "NotFound"
	case CodePermissionDenied:
		return "PermissionDenied"
	default:
		return "Unknown"
	}
}

// --- Metadata (имитация grpc metadata) ---

type metadataKey struct{}

// WithMetadata добавляет пары key-value в context.
func WithMetadata(ctx context.Context, kv ...string) context.Context {
	m := make(map[string]string)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i]] = kv[i+1]
	}
	return context.WithValue(ctx, metadataKey{}, m)
}

// GetMetadata извлекает значение из metadata context.
func GetMetadata(ctx context.Context, key string) (string, bool) {
	m, ok := ctx.Value(metadataKey{}).(map[string]string)
	if !ok {
		return "", false
	}
	v, found := m[key]
	return v, found
}

// --- Счётчики метрик ---

type metricsStore struct {
	calls  map[string]int
	errors map[string]int
}

var metrics = &metricsStore{
	calls:  make(map[string]int),
	errors: make(map[string]int),
}

// --- Интерцепторы ---

// TODO: реализуй LoggingInterceptor() UnaryInterceptor
// Шаги:
//   1. Запомни время начала: start := time.Now()
//   2. Вызови handler и получи результат + ошибку
//   3. Вычисли длительность: elapsed := time.Since(start)
//   4. Определи статус: если err == nil → "OK", иначе codeName(se.Code)
//   5. Напечатай: [log] <method> <elapsed.Round(100μs)> <status>

func LoggingInterceptor() UnaryInterceptor {
	return func(ctx context.Context, method string, req any, handler UnaryHandler) (any, error) {
		// TODO: implement
		start := time.Now()

		resp, err := handler(ctx, req)

		elapsed := time.Since(start) // так посчитали время на отработку

		status := "OK"

		if err != nil {
			if se, ok := err.(*StatusError); ok {
				status = codeName(se.Code)
			} else {
				status = "Unknown"
			}
		}

		fmt.Printf("[log] %s %s %s /n", method, elapsed.Round(100*time.Microsecond), status)

		return resp, err
	}
}

// TODO: реализуй AuthInterceptor(serviceToken string) UnaryInterceptor
// Шаги:
//   1. Извлеки "authorization" из metadata: GetMetadata(ctx, "authorization")
//   2. Если значение не равно serviceToken → верни StatusError{CodePermissionDenied, "invalid service token"}
//   3. Иначе — вызови handler и верни результат

func AuthInterceptor(serviceToken string) UnaryInterceptor {

	return func(
		ctx context.Context,
		method string,
		req any,
		handler UnaryHandler,
	) (any, error) {

		token, ok := GetMetadata( // пытаемся достать мето данные по ключу
			ctx,
			"authorization",
		)

		if !ok || token != serviceToken { // так если вообще нет или токен не правильный то возвращем ошибку
			return nil, &StatusError{
				Code:    CodePermissionDenied,
				Message: "invalid service token",
			}
		}

		return handler(ctx, req)
	}
}

// TODO: реализуй MetricsInterceptor() UnaryInterceptor
// Шаги:
//   1. Инкрементируй metrics.calls[method]
//   2. Вызови handler
//   3. Если err != nil — инкрементируй metrics.errors[method]
//   4. Напечатай: [metrics] <method> calls=<N>
//   5. Верни результат

func MetricsInterceptor() UnaryInterceptor {
	return func(ctx context.Context, method string, req any, handler UnaryHandler) (any, error) {

		metrics.calls[method]++ // счетчик вызщова метода

		resp, err := handler(ctx, req)

		if err != nil {
			metrics.errors[method]++
		}

		fmt.Printf(
			"[metrics] %s calls=%d\n",
			method,
			metrics.calls[method],
		)

		return resp, err
	}
}

// --- Chain ---

// TODO: реализуй ChainInterceptors(interceptors ...UnaryInterceptor) UnaryInterceptor
// Цепочка выполняется слева направо: первый интерцептор вызывает второй и т.д.
// Подсказка: рекурсивно или через замыкание.

func ChainInterceptors(interceptors ...UnaryInterceptor) UnaryInterceptor {
	return func(ctx context.Context, method string, req any, handler UnaryHandler) (any, error) {
		// TODO: implement
		chainedHandler := handler

		for i := len(interceptors) - 1; i >= 0; i-- {

			interceptor := interceptors[i]
			next := chainedHandler

			chainedHandler = func(
				ctx context.Context,
				req any,
			) (any, error) {

				return interceptor(
					ctx,
					method,
					req,
					next,
				)
			}
		}

		return chainedHandler(ctx, req)
	}
}

// --- Сервис ---

type IssueTokenRequest struct{ UserID string }
type IssueTokenResponse struct{ TokenID, UserID string }
type ValidateTokenRequest struct{ TokenID string }
type ValidateTokenResponse struct {
	UserID string
	Valid  bool
}

type tokenStore struct {
	tokens map[string]string
	nextID int
}

func newTokenStore() *tokenStore {
	return &tokenStore{tokens: make(map[string]string), nextID: 1}
}

func (s *tokenStore) issueToken(_ context.Context, req any) (any, error) {
	r := req.(*IssueTokenRequest)
	if r.UserID == "" {
		return nil, &StatusError{Code: CodeInvalidArgument, Message: "user_id is required"}
	}
	id := fmt.Sprintf("tok-%d", s.nextID)
	s.nextID++
	s.tokens[id] = r.UserID
	return &IssueTokenResponse{TokenID: id, UserID: r.UserID}, nil
}

func (s *tokenStore) validateToken(_ context.Context, req any) (any, error) {
	r := req.(*ValidateTokenRequest)
	uid, ok := s.tokens[r.TokenID]
	if !ok {
		return nil, &StatusError{Code: CodeNotFound, Message: "token not found"}
	}
	return &ValidateTokenResponse{UserID: uid, Valid: true}, nil
}

func call(ctx context.Context, method string, req any, handler UnaryHandler, chain UnaryInterceptor) (any, error) {
	return chain(ctx, method, req, handler)
}

func main() {
	store := newTokenStore()
	chain := ChainInterceptors(
		MetricsInterceptor(),
		LoggingInterceptor(),
		AuthInterceptor("secret-service-token"),
	)

	// Запрос с правильным токеном сервиса
	ctx := WithMetadata(context.Background(), "authorization", "secret-service-token")

	resp, err := call(ctx, "/token.TokenService/IssueToken", &IssueTokenRequest{UserID: "u-1"}, store.issueToken, chain)
	if err != nil {
		fmt.Println("IssueToken error:", err)
		return
	}
	issued := resp.(*IssueTokenResponse)
	fmt.Println("Issued:", issued.TokenID)

	resp, err = call(ctx, "/token.TokenService/ValidateToken", &ValidateTokenRequest{TokenID: issued.TokenID}, store.validateToken, chain)
	if err != nil {
		fmt.Println("ValidateToken error:", err)
	} else {
		fmt.Println("Validated:", resp.(*ValidateTokenResponse).UserID)
	}

	// Запрос без токена сервиса — ожидаем PermissionDenied
	ctxNoAuth := context.Background()
	_, err = call(ctxNoAuth, "/token.TokenService/ValidateToken", &ValidateTokenRequest{TokenID: issued.TokenID}, store.validateToken, chain)
	if err != nil {
		fmt.Println("No auth error:", err)
	}
}
