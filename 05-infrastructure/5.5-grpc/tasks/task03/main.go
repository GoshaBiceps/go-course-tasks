// Задание: Реализация gRPC-сервера
//
// Реализуй полноценный gRPC-сервер для TokenService.
// В реальном проекте контракт описывается в .proto файле;
// здесь используем Go-интерфейсы как эквивалент.
//
// Методы сервиса:
//   IssueToken(ctx, req *IssueTokenRequest) (*IssueTokenResponse, error)
//   ValidateToken(ctx, req *ValidateTokenRequest) (*ValidateTokenResponse, error)
//   RevokeToken(ctx, req *RevokeTokenRequest) (*RevokeTokenResponse, error)
//
// Требования:
//   - Хранилище токенов in-memory (map)
//   - Генерация ID: "tok-1", "tok-2", ...
//   - Пустой UserID → ошибка (gRPC InvalidArgument)
//   - Несуществующий / отозванный токен → ошибка (gRPC NotFound / PermissionDenied)
//   - Имитируй gRPC-статусы через пакет google.golang.org/grpc/codes + status
//
// Ожидаемый результат:
//   IssueToken("u-1")   → {TokenID:"tok-1", UserID:"u-1"}
//   ValidateToken(tok-1) → {UserID:"u-1", Valid:true}
//   RevokeToken(tok-1)   → {Revoked:true}
//   ValidateToken(tok-1) → error: token revoked
//
// Зависимости:
//   go get google.golang.org/grpc

package main

import (
	"context"
	"errors"
	"fmt"
)

// gRPC-коды статусов (имитация без импорта grpc/codes)
type StatusCode int

const ( // коды
	CodeOK               StatusCode = 0
	CodeInvalidArgument  StatusCode = 3
	CodeNotFound         StatusCode = 5
	CodePermissionDenied StatusCode = 7
)

type StatusError struct { // ошибка с статус кодом
	Code    StatusCode
	Message string
}

func (e *StatusError) Error() string { // текст ошибки
	return fmt.Sprintf("rpc error: code = %d desc = %s", e.Code, e.Message)
}

func statusError(code StatusCode, msg string) error { //  возвращение нашей кастомной ошибки
	return &StatusError{Code: code, Message: msg}
}

// --- Типы запросов и ответов ---

type IssueTokenRequest struct { // запрос на выдачу токена
	UserID string
}

type IssueTokenResponse struct { // ответ после выдачи токена
	TokenID string
	UserID  string
}

type ValidateTokenRequest struct { //  запрос на проверку токена
	TokenID string
}

type ValidateTokenResponse struct { //  ну соответственно ответ проверки
	UserID string
	Valid  bool
}

type RevokeTokenRequest struct { //  так запрос на отзыв токена
	TokenID string
}

type RevokeTokenResponse struct { // ну и ответ токен отозван
	Revoked bool
}

// --- Интерфейс сервиса ---

type TokenServiceServer interface {
	IssueToken(ctx context.Context, req *IssueTokenRequest) (*IssueTokenResponse, error)
	ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error)
	RevokeToken(ctx context.Context, req *RevokeTokenRequest) (*RevokeTokenResponse, error)
}

// --- Реализация ---

type tokenRecord struct {
	tokenID string
	userID  string
	revoked bool
}

type tokenServiceImpl struct {
	tokens map[string]*tokenRecord
	nextID int
}

func NewTokenService() TokenServiceServer {
	return &tokenServiceImpl{
		tokens: make(map[string]*tokenRecord),
		nextID: 1,
	}
}

// TODO: реализуй IssueToken:
//   1. Проверь UserID != "" → иначе statusError(CodeInvalidArgument, "user_id is required")
//   2. Сгенерируй tokenID: fmt.Sprintf("tok-%d", s.nextID)
//   3. Сохрани в s.tokens
//   4. Верни IssueTokenResponse

func (s *tokenServiceImpl) IssueToken(ctx context.Context, req *IssueTokenRequest) (*IssueTokenResponse, error) {

	if req.UserID == "" {
		return nil, statusError(CodeInvalidArgument, "user_id is reqyured")
	}

	tokenID := fmt.Sprintf("tok-%d", s.nextID)

	record := &tokenRecord{
		tokenID: tokenID,
		userID:  req.UserID,
		revoked: false,
	}

	s.tokens[record.tokenID] = record
	s.nextID++

	return &IssueTokenResponse{
		TokenID: record.tokenID,
		UserID:  record.userID,
	}, nil
}

// TODO: реализуй ValidateToken:
//   1. Найди токен в s.tokens → если нет: statusError(CodeNotFound, "token not found")
//   2. Если revoked: statusError(CodePermissionDenied, "token revoked")
//   3. Верни ValidateTokenResponse{UserID: ..., Valid: true}

func (s *tokenServiceImpl) ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	// TODO: implement
	record, ok := s.tokens[req.TokenID]
	if !ok {
		return nil, statusError(CodeNotFound, "token not found")
	}

	if record.revoked == true {
		return nil, statusError(CodePermissionDenied, "token revoked")
	}

	return &ValidateTokenResponse{
		UserID: record.userID,
		Valid:  true,
	}, nil
}

// TODO: реализуй RevokeToken:
//   1. Найди токен → если нет: statusError(CodeNotFound, "token not found")
//   2. Установи revoked = true
//   3. Верни RevokeTokenResponse{Revoked: true}

func (s *tokenServiceImpl) RevokeToken(ctx context.Context, req *RevokeTokenRequest) (*RevokeTokenResponse, error) {

	token, ok := s.tokens[req.TokenID]
	if !ok {
		return nil, statusError(CodeNotFound, "token not found")
	}

	token.revoked = true

	// TODO: implement
	return &RevokeTokenResponse{
		Revoked: token.revoked,
	}, nil
}

func main() {
	svc := NewTokenService()
	ctx := context.Background()

	// Issue token
	issued, err := svc.IssueToken(ctx, &IssueTokenRequest{UserID: "u-1"})
	if err != nil {
		fmt.Println("IssueToken error:", err)
		return
	}
	fmt.Printf("IssueToken: %+v\n", issued)

	// Validate token
	validated, err := svc.ValidateToken(ctx, &ValidateTokenRequest{TokenID: issued.TokenID})
	if err != nil {
		fmt.Println("ValidateToken error:", err)
	} else {
		fmt.Printf("ValidateToken: %+v\n", validated)
	}

	// Revoke token
	revoked, err := svc.RevokeToken(ctx, &RevokeTokenRequest{TokenID: issued.TokenID})
	if err != nil {
		fmt.Println("RevokeToken error:", err)
	} else {
		fmt.Printf("RevokeToken: %+v\n", revoked)
	}

	// Validate revoked token — ожидаем ошибку
	_, err = svc.ValidateToken(ctx, &ValidateTokenRequest{TokenID: issued.TokenID})
	var se *StatusError
	if errors.As(err, &se) && se.Code == CodePermissionDenied {
		fmt.Println("ValidateToken (revoked): permission denied (expected)")
	}

	// Issue with empty UserID — ожидаем ошибку
	_, err = svc.IssueToken(ctx, &IssueTokenRequest{UserID: ""})
	if errors.As(err, &se) && se.Code == CodeInvalidArgument {
		fmt.Println("IssueToken (empty user): invalid argument (expected)")
	}
}
