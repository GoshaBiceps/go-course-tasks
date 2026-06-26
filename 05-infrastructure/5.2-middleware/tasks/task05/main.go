// Задание: JWT/Paseto адаптеры
//
// Сделай две реализации TokenVerifier.
// HTTP-слой должен работать без изменений при смене реализации.
//
// Ожидаемый результат:
//   using JWTVerifier:
//   verified: {UserID:user-jwt Role:user}
//   error: jwt: unknown token
//
//   using PasetoVerifier:
//   verified: {UserID:user-paseto Role:admin}
//   error: paseto: unknown token

package main

import (
	"errors"
	"fmt"
)

type Claims struct {
	UserID string
	Role   string
}

type TokenVerifier interface {
	Verify(token string) (Claims, error)
}

// TODO: реализуй JWTVerifier (мок, без реальной криптографии)
// Если token == "jwt-valid" — вернуть Claims{UserID: "user-jwt", Role: "user"}
// Иначе — вернуть ошибку "jwt: unknown token"

type JWTVerifier struct{}

func (v *JWTVerifier) Verify(token string) (Claims, error) {
	if token != "jwt-valid" {
		return Claims{}, errors.New("jwt: unknown token")
	}

	return Claims{UserID: "user-jwt", Role: "user"}, nil
}

// TODO: реализуй PasetoVerifier (мок, без реальной криптографии)
// Если token == "paseto-valid" — вернуть Claims{UserID: "user-paseto", Role: "admin"}
// Иначе — вернуть ошибку "paseto: unknown token"

type PasetoVerifier struct{}

func (v *PasetoVerifier) Verify(token string) (Claims, error) {
	if token != "paseto-valid" {
		return Claims{}, errors.New("pseto: unknown token")
	}
	return Claims{UserID: "user-paseto", Role: "admin"}, nil
}

func runVerification(name string, verifier TokenVerifier, tokens []string) {
	fmt.Printf("using %s:\n", name)
	for _, t := range tokens {
		claims, err := verifier.Verify(t)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Printf("verified: %+v\n", claims)
		}
	}
	fmt.Println()
}

func main() {

	runVerification("JWTVerifier", &JWTVerifier{}, []string{"jwt-valid", "bad-token"})
	runVerification("PasetoVerifier", &PasetoVerifier{}, []string{"paseto-valid", "bad-token"})
}
