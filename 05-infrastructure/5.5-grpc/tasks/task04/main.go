// Задание: Реализация gRPC-клиента
//
// Реализуй клиент для TokenService.
// В реальном проекте клиент генерируется из .proto (protoc --go-grpc_out=.);
// здесь используем интерфейс и in-process "сервер" вместо сетевого подключения.
//
// Задачи:
//   1. Реализуй TokenServiceClient с методами IssueToken / ValidateToken / RevokeToken
//   2. Клиент принимает адрес сервера и "подключается" (imitate dial)
//   3. Обрабатывай gRPC-статусы: NotFound, PermissionDenied, InvalidArgument
//
// Ожидаемый результат:
//   Dialing token-service at localhost:50051...
//   Issued token: tok-1 for user u-1
//   Token tok-1 is valid for user u-1
//   Token tok-1 revoked
//   ValidateToken after revoke: rpc error: code = 7 desc = token revoked
//
// В реальном проекте:
//   conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
//   client := pb.NewTokenServiceClient(conn)

package main

import (
	"context"
	"errors"
	"fmt"
)

// --- Типы (дублируем из task03 для автономности задачи) ---

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

type IssueTokenRequest struct{ UserID string }
type IssueTokenResponse struct{ TokenID, UserID string }
type ValidateTokenRequest struct{ TokenID string }
type ValidateTokenResponse struct {
	UserID string
	Valid  bool
}
type RevokeTokenRequest struct{ TokenID string }
type RevokeTokenResponse struct{ Revoked bool }

// --- Интерфейс клиента ---

type TokenServiceClient interface {
	IssueToken(ctx context.Context, req *IssueTokenRequest) (*IssueTokenResponse, error)
	ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error)
	RevokeToken(ctx context.Context, req *RevokeTokenRequest) (*RevokeTokenResponse, error)
}

// --- In-process "сервер" (заменяет реальный gRPC-сервер) ---

type inProcessServer struct {
	tokens map[string]*tokenRecord
	nextID int
}

type tokenRecord struct {
	tokenID string
	userID  string
	revoked bool
}

func newInProcessServer() *inProcessServer {
	return &inProcessServer{tokens: make(map[string]*tokenRecord), nextID: 1}
}

func (s *inProcessServer) IssueToken(_ context.Context, req *IssueTokenRequest) (*IssueTokenResponse, error) {
	if req.UserID == "" {
		return nil, &StatusError{Code: CodeInvalidArgument, Message: "user_id is required"}
	}
	id := fmt.Sprintf("tok-%d", s.nextID)
	s.nextID++
	s.tokens[id] = &tokenRecord{tokenID: id, userID: req.UserID}
	return &IssueTokenResponse{TokenID: id, UserID: req.UserID}, nil
}

func (s *inProcessServer) ValidateToken(_ context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	t, ok := s.tokens[req.TokenID]
	if !ok {
		return nil, &StatusError{Code: CodeNotFound, Message: "token not found"}
	}
	if t.revoked {
		return nil, &StatusError{Code: CodePermissionDenied, Message: "token revoked"}
	}
	return &ValidateTokenResponse{UserID: t.userID, Valid: true}, nil
}

func (s *inProcessServer) RevokeToken(_ context.Context, req *RevokeTokenRequest) (*RevokeTokenResponse, error) {
	t, ok := s.tokens[req.TokenID]
	if !ok {
		return nil, &StatusError{Code: CodeNotFound, Message: "token not found"}
	}
	t.revoked = true
	return &RevokeTokenResponse{Revoked: true}, nil
}

// --- Клиент ---

// TODO: реализуй tokenServiceClient, который обёртывает in-process сервер.
// В реальном проекте здесь был бы grpc.ClientConn.

type tokenServiceClient struct {
	server *inProcessServer // заменяет реальное соединение
}

// TODO: реализуй dial(addr string) (TokenServiceClient, error)
//   Напечатай "Dialing token-service at <addr>..."
//   Верни tokenServiceClient, обёртывающий inProcessServer

func dial(addr string) (TokenServiceClient, error) {

	fmt.Printf("Dailing token-service at %s...\n", addr)

	server := newInProcessServer()

	client := &tokenServiceClient{
		server: server,
	}

	return client, nil
}

// TODO: реализуй методы tokenServiceClient — делегируй вызовы на s.server
func (c *tokenServiceClient) IssueToken(
	ctx context.Context,
	req *IssueTokenRequest,
) (*IssueTokenResponse, error) {

	return c.server.IssueToken(ctx, req)
}

func (c *tokenServiceClient) ValidateToken(
	ctx context.Context,
	req *ValidateTokenRequest,
) (*ValidateTokenResponse, error) {

	return c.server.ValidateToken(ctx, req)
}

func (c *tokenServiceClient) RevokeToken(
	ctx context.Context,
	req *RevokeTokenRequest,
) (*RevokeTokenResponse, error) {

	return c.server.RevokeToken(ctx, req)
}

// TODO: реализуй handleError(err error) — проверяй StatusError.Code и печатай описание:
//   CodeNotFound        → "not found: <msg>"
//   CodePermissionDenied → "permission denied: <msg>"
//   CodeInvalidArgument → "invalid argument: <msg>"
//   иначе               → "unexpected error: <err>"

func handleError(err error) {

	var se *StatusError

	if errors.As(err, &se) {

		switch se.Code {

		case CodeNotFound:
			fmt.Println("not found:", se.Message)

		case CodePermissionDenied:
			fmt.Println("permission denied:", se.Message)

		case CodeInvalidArgument:
			fmt.Println("invalid argument:", se.Message)

		default:
			fmt.Println("unexpected error:", err)
		}

		return
	}

	fmt.Println("unexpected error:", err)
}

func main() {
	ctx := context.Background()

	client, err := dial("localhost:50051")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}

	// Issue token
	issued, err := client.IssueToken(ctx, &IssueTokenRequest{UserID: "u-1"})
	if err != nil {
		handleError(err)
		return
	}
	fmt.Printf("Issued token: %s for user %s\n", issued.TokenID, issued.UserID)

	// Validate token
	validated, err := client.ValidateToken(ctx, &ValidateTokenRequest{TokenID: issued.TokenID})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("Token %s is valid for user %s\n", issued.TokenID, validated.UserID)
	}

	// Revoke token
	_, err = client.RevokeToken(ctx, &RevokeTokenRequest{TokenID: issued.TokenID})
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("Token %s revoked\n", issued.TokenID)
	}

	// Validate after revoke
	_, err = client.ValidateToken(ctx, &ValidateTokenRequest{TokenID: issued.TokenID})
	if err != nil {
		fmt.Println("ValidateToken after revoke:", err)
	}
}
