// Задание: Bearer token parser
//
// Реализуй функцию extractBearerToken, которая корректно разбирает
// заголовок Authorization.
//
// Ожидаемый результат:
//   extracted: my-secret-token
//   error: authorization header is empty
//   error: invalid authorization header format
//   error: token is empty

package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyHeader   = errors.New("authorization header is empty")
	ErrInvalidFormat = errors.New("invalid authorization header format")
	ErrEmptyToken    = errors.New("token is empty")
)

// TODO: реализуй функцию extractBearerToken(header string) (string, error)
// Логика:
//   1. Если header == "" — вернуть ErrEmptyHeader
//   2. Если не начинается с "Bearer " — вернуть ErrInvalidFormat
//   3. Достать токен: часть после "Bearer "
//   4. Если токен == "" — вернуть ErrEmptyToken
//   5. Вернуть токен и nil
// Подсказка: используй strings.CutPrefix или strings.HasPrefix

func extractBearerToken(header string) (string, error) {
	const prefix = "Bearer "

	if header == "" { //  если header пришел пустой
		return "", ErrEmptyHeader
	}

	if !strings.HasPrefix(header, prefix) {
		return "", ErrInvalidFormat // не валидный формат
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix)) // достали чистый токен

	if token == "" {
		return "", ErrEmptyToken // пустой токен
	}

	return token, nil
}

func main() {
	cases := []struct {
		header string
	}{
		{"Bearer my-secret-token"},
		{""},
		{"Token my-secret-token"},
		{"Bearer "},
	}

	for _, c := range cases {
		token, err := extractBearerToken(c.header)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("extracted:", token)
		}
	}
}
