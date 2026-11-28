package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/service"
)

type MockUserService struct {
	registerCalled bool
	registerErr    error
	loginCalled    bool
	loginToken     string
	loginErr       error
}

func (m *MockUserService) Register(ctx context.Context, email string, password string) error {
	m.registerCalled = true
	return m.registerErr
}

func (m *MockUserService) Login(ctx context.Context, email string, password string) (string, error) {
	m.loginCalled = true
	return m.loginToken, m.loginErr
}

func TestRegister_InvalidJSON(t *testing.T) {
	h := NewUserHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("{invalid json"))

	rec := httptest.NewRecorder() // ответ

	h.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

}

func TestRegister_InvalidData(t *testing.T) {
	mock := &MockUserService{
		registerErr: service.ErrInvalidRegistrationData,
	}

	h := NewUserHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"email":"invalid-email","password":"12345678"}`))

	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "invalid registration data") {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}

	if !mock.registerCalled {
		t.Fatal("expected Register to be called")
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mock := &MockUserService{
		registerErr: service.ErrEmailAlreadyExists,
	}

	h := NewUserHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"email":"test@mail.com","password":"12345678"}`))

	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "email already registered") {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}

	if !mock.registerCalled {
		t.Fatal("expected Register to be called")
	}
}

func TestRegister_InternalError(t *testing.T) {
	mock := &MockUserService{
		registerErr: errors.New("database unavailable"),
	}

	h := NewUserHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"email":"test@mail.com","password":"12345678"}`))

	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "internal server error") {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}

	if !mock.registerCalled {
		t.Fatal("expected Register to be called")
	}
}

func TestRegister_Success(t *testing.T) {
	mock := &MockUserService{}

	h := NewUserHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/register",
		strings.NewReader(`{"email":"test@mail.com","password":"12345678"}`))

	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "user registered successfully") {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}

	if !mock.registerCalled {
		t.Fatal("expected Register to be called")
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	h := NewUserHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("{invalid json"))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "invalid JSON") {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mock := &MockUserService{
		loginErr: service.ErrInvalidCredentials,
	}

	h := NewUserHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"test@mail.com","password":"wrongpassword"}`))

	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "invalid email or password") {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}

	if !mock.loginCalled {
		t.Fatal("expected Login to be called")
	}
}
func TestLogin_InternalError(t *testing.T) {
	mock := &MockUserService{
		loginErr: errors.New("database unavailable"),
	}

	h := NewUserHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"test@mail.com","password":"12345678"}`))

	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "internal server error") {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}

	if !mock.loginCalled {
		t.Fatal("expected Login to be called")
	}
}

func TestLogin_Success(t *testing.T) {
	mock := &MockUserService{
		loginToken: "fake-jwt-token",
	}

	h := NewUserHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"test@mail.com","password":"12345678"}`))

	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected Content-Type: %s", rec.Header().Get("Content-Type"))
	}

	var body map[string]string

	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["token"] != "fake-jwt-token" {
		t.Fatalf("unexpected token: %s", body["token"])
	}

	if !mock.loginCalled {
		t.Fatal("expected Login to be called")
	}
}
