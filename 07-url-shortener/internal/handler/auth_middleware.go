package handler

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header is required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") { // проверяем чтобы строка начиналась с Bearer
			http.Error(
				w, "invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ") // так удаляем начало строки , чтобы остался один токен

		if tokenString == "" {
			http.Error(w, "token is required", http.StatusUnauthorized) // проверка, если убрали начало записи . есть сам токен или нет
			return
		}

		claims := &jwt.RegisteredClaims{} // сюда мы запишем данные токена после разбора

		parsedToken, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(_ *jwt.Token) (any, error) { // функция разбирате и проверяет данные с токеном
				return []byte(jwtSecret), nil // возвращем наш секретный ключ в байтах
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), // разрешаем принимать jwt только подписанные алкгортмом hs256
		)
		if err != nil || !parsedToken.Valid {

			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		next(w, r) // вызываем наш хендлер
	}
}
