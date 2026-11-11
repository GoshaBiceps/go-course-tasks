package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware_NoToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/shorten", nil) // наш искуственый запрос
	rec := httptest.NewRecorder()                                // место куда запишеться Http ответ

	next := func(w http.ResponseWriter, r *http.Request) { //
		t.Error("next must not be collded")
	}

	handler := AuthMiddleware("test-secret", next)

	handler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secret := "test-secret" // тестовый секретный ключ

	token := jwt.NewWithClaims( //  создаем токен
		jwt.SigningMethodHS256, // так алгоритм подписи
		jwt.RegisteredClaims{}, // пустой набор
	)

	tokenString, err := token.SignedString([]byte(secret)) //подписываем наш токен (перевод в байты потому тчо метод ожидает именно этот тип)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err) // сообщаем что тест провален и останавливаем его
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", nil)

	req.Header.Set("Authorization", "Bearer "+tokenString) // так добаляем JWT в заголовок запроса.

	rec := httptest.NewRecorder() // сюда записываем наш ответ клиенту

	nextCalled := false // вызов следующего обработчика

	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusNoContent)
	}

	handler := AuthMiddleware(secret, next)

	handler(rec, req) // запускаем проверку авторизации .

	if !nextCalled {
		t.Error("next handler was not called")
	}

	if rec.Code != http.StatusNoContent {
		t.Errorf(
			"expected status 204, got %d",
			rec.Code,
		)
	}
}

func TestAuthMIddleware_InvalidSignature(t *testing.T) {
	secret := "test-secret"

	token := jwt.NewWithClaims( //  создаем токен
		jwt.SigningMethodHS256, // так алгоритм подписи
		jwt.RegisteredClaims{}, // пустой набор
	)

	tokenString, err := token.SignedString([]byte(secret)) // подписываем токен
	if err != nil {
		t.Fatalf("failed to sign token: %v", err) // сообщаем что тест провален и останавливаем его
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", nil)

	req.Header.Set("Authorization", "Bearer "+tokenString)

	rec := httptest.NewRecorder()

	nextCalled := false // вызов следующего обработчика

	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusNoContent)
	}

	handler := AuthMiddleware("wrong-secret", next)

	handler(rec, req) // запускаем проверку авторизации .

	if nextCalled {
		t.Error("next handler was not called")
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf(
			"expected status 204, got %d",
			rec.Code,
		)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/shorten", nil) // наш искуственый запрос
	req.Header.Set("Authorization", "abc123")
	rec := httptest.NewRecorder() // место куда запишеться Http ответ

	next := func(w http.ResponseWriter, r *http.Request) { //
		t.Error("next must not be collded")
	}

	handler := AuthMiddleware("test-secret", next)

	handler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}
