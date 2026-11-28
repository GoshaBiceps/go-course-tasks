package service

import (
	"context"
	"errors"
	"testing"
	"url-shortener/internal/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type MockUserStore struct {
	createUserCalled bool                      // вызывали ли метод
	createUserParams database.CreateUserParams // данные юзера email и хеш
	createUserErr    error                     // ошибка которую вернем при создании юзера
	getUserCalled    bool                      // вызывали ли метод
	getUserResult    database.User             //пользватель , которого вернем при вызове CreateUser
	getUserErr       error                     // ошибка при вызове GetUserByEmail
}

func (m *MockUserStore) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.CreateUserRow, error) {
	m.createUserCalled = true
	m.createUserParams = arg
	return database.CreateUserRow{}, m.createUserErr
}

func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	m.getUserCalled = true
	return m.getUserResult, m.getUserErr
}

func TestRegister_EmptyEmail(t *testing.T) {
	mock := &MockUserStore{}
	svc := NewUserService(mock, "test-secret")

	err := svc.Register(
		context.Background(),
		"",
		"12345678",
	)

	if !errors.Is(err, ErrInvalidRegistrationData) {
		t.Fatalf(
			"expected ErrInvalidRegistrationData, got %v",
			err,
		)
	}

	if mock.createUserCalled {
		t.Fatal("CreateUser should not be called")
	}
}

func TestRegister_EmptyPassword(t *testing.T) {
	mock := &MockUserStore{}
	svc := NewUserService(mock, "test-secret")

	err := svc.Register(
		context.Background(),
		"test@gmail.com",
		"",
	)

	if !errors.Is(err, ErrInvalidRegistrationData) {
		t.Fatalf(
			"expected ErrInvalidRegistrationData, got %v",
			err,
		)
	}

	if mock.createUserCalled {
		t.Fatal("CreateUser should not be called")
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	mock := &MockUserStore{}
	svc := NewUserService(mock, "test-secret")

	err := svc.Register(
		context.Background(),
		"test@gmail.com",
		"1234567",
	)

	if !errors.Is(err, ErrInvalidRegistrationData) {
		t.Fatalf("expected ErrInvalidRegistrationData, got %v", err)
	}

	if mock.createUserCalled {
		t.Fatal("CreateUser should not be called")
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	mock := &MockUserStore{}
	svc := NewUserService(mock, "test-secret")

	err := svc.Register(
		context.Background(),
		"testgmail.com",
		"12345678",
	)

	if !errors.Is(err, ErrInvalidRegistrationData) {
		t.Fatalf("expected ErrInvalidRegistrationData, got %v", err)
	}

	if mock.createUserCalled {
		t.Fatal("CreateUser should not be called")
	}
}

func TestRegister_EmailWithDisplayName(t *testing.T) {
	mock := &MockUserStore{}
	svc := NewUserService(mock, "test-secret")

	err := svc.Register(
		context.Background(),
		"Ivan <ivan@gmail.com>",
		"12345678",
	)

	if !errors.Is(err, ErrInvalidRegistrationData) {
		t.Fatalf("expected ErrInvalidRegistrationData, got %v", err)
	}

	if mock.createUserCalled {
		t.Fatal("CreateUser should not be called")
	}
}

func TestRegister_Success(t *testing.T) {
	mock := &MockUserStore{}
	svc := NewUserService(mock, "test-secret")

	email := "test@gmail.com"
	password := "12345678"

	err := svc.Register(context.Background(), email, password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !mock.createUserCalled {
		t.Fatal("expected CreateUser to be called")
	}

	// Проверяем, что в репозиторий передан правильный email
	if mock.createUserParams.Email != email {
		t.Fatalf("expected email %s, got %s", email, mock.createUserParams.Email)
	}

	// Проверяем, что пароль не сохранён в открытом виде
	if mock.createUserParams.PasswordHash == password {
		t.Fatal("password must be hashed")
	}

	// Проверяем, что сохранённый хеш соответствует паролю
	err = bcrypt.CompareHashAndPassword(
		[]byte(mock.createUserParams.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		t.Fatalf("password hash verification failed: %v", err)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mock := &MockUserStore{
		createUserErr: &pgconn.PgError{
			Code:           "23505",
			ConstraintName: "users_email_key",
		},
	}
	svc := NewUserService(mock, "test-secret")

	err := svc.Register(context.Background(), "test@gmail.com", "12345678")
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestRegister_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database unavailable")

	mock := &MockUserStore{
		createUserErr: expectedErr,
	}

	svc := NewUserService(mock, "test-secret")

	err := svc.Register(
		context.Background(),
		"test@gmail.com",
		"12345678",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if !mock.createUserCalled {
		t.Fatal("expected CreateUser to be called")
	}
}

func TestLogin_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database unavailable")

	mock := &MockUserStore{
		getUserErr: expectedErr,
	}

	svc := NewUserService(mock, "test-secret")

	token, err := svc.Login(
		context.Background(),
		"test@gmail.com",
		"12345678",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if token != "" {
		t.Fatalf("expected empty token, got %s", token)
	}

	if !mock.getUserCalled {
		t.Fatal("expected GetUserByEmail to be called")
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	password := "12345678"

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.MinCost,
	)
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}

	mock := &MockUserStore{
		getUserResult: database.User{
			ID:           42,
			PasswordHash: string(hash),
		},
	}

	svc := NewUserService(mock, "test-secret")

	token, err := svc.Login(
		context.Background(),
		"test@gmail.com",
		"wrongpassword",
	)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}

	if token != "" {
		t.Fatalf("expected empty token, got %s", token)
	}

	if !mock.getUserCalled {
		t.Fatal("expected GetUserByEmail to be called")
	}
}

func TestLogin_InvalidHash(t *testing.T) {
	mock := &MockUserStore{
		getUserResult: database.User{
			ID:           42,
			PasswordHash: "invalid-hash",
		},
	}

	svc := NewUserService(mock, "test-secret")

	token, err := svc.Login(
		context.Background(),
		"test@gmail.com",
		"12345678",
	)

	if !errors.Is(err, bcrypt.ErrHashTooShort) {
		t.Fatalf("expected bcrypt.ErrHashTooShort, got %v", err)
	}

	if token != "" {
		t.Fatalf("expected empty token, got %s", token)
	}
}

func TestLogin_Success(t *testing.T) {
	password := "12345678"

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.MinCost,
	)
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}

	mock := &MockUserStore{
		getUserResult: database.User{
			ID:           42,
			PasswordHash: string(hash),
		},
	}

	svc := NewUserService(mock, "test-secret")

	tokenString, err := svc.Login(
		context.Background(),
		"test@gmail.com",
		password,
	)

	if err != nil {
		t.Fatalf("expected successful login, got %v", err)
	}

	if tokenString == "" {
		t.Fatal("expected non-empty JWT token")
	}

	if !mock.getUserCalled {
		t.Fatal("expected GetUserByEmail to be called")
	}

	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte("test-secret"), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil || !token.Valid {
		t.Fatalf("expected valid JWT, got %v", err)
	}

	if claims.Subject != "42" {
		t.Fatalf("expected user ID 42, got %s", claims.Subject)
	}

	if claims.ExpiresAt == nil {
		t.Fatal("expected JWT expiration time")
	}
}
