package service

import (
	"context"
	"errors"
	"testing"
	"url-shortener/internal/database"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type MockUserStore struct {
	createUserCalled bool
	createUserParams database.CreateUserParams
	createUserErr    error
}

func (m *MockUserStore) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.CreateUserRow, error) {
	m.createUserCalled = true
	m.createUserParams = arg
	return database.CreateUserRow{}, m.createUserErr
}

func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	return database.User{}, nil
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
