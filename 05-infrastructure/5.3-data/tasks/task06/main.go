// Задание: Интеграционная мини-задача — data-слой Token Service
//
// Спроектируй и реализуй data-слой для Token Service.
//
// Таблицы:
//   users         (id, email, password_hash, created_at)
//   refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at)
//
// Методы репозитория:
//   UserRepository:
//     CreateUser(ctx, email, passwordHash string) (User, error)
//     GetUserByEmail(ctx, email string) (User, error)
//
//   TokenRepository:
//     CreateRefreshToken(ctx, userID int64, tokenHash string, expiresAt time.Time) (RefreshToken, error)
//     GetActiveToken(ctx, tokenHash string) (RefreshToken, error)  — только не отозванные и не истёкшие
//     RevokeToken(ctx, tokenHash string) error
//
// Требования:
//   - PostgreSQL + pgxpool
//   - type-safe слой через sqlc (или ручной pgx)
//   - миграции в sql/migrations/
//   - транзакция для сценария "создание пользователя + выдача refresh-token"
//   - краткий README с описанием схемы
//
// Структура директории:
//   task06/
//     main.go              ← этот файл
//     go.mod
//     sql/
//       migrations/
//         001_users.up.sql
//         001_users.down.sql
//         002_tokens.up.sql
//         002_tokens.down.sql
//       queries/
//         users.sql
//         tokens.sql
//     db/                  ← сгенерированный sqlc-код (или ручные репозитории)
//     README.md

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cource/data-task06/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenRevoked  = errors.New("token revoked")
	ErrTokenExpired  = errors.New("token expired")
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type userRepository struct { // для рабты  в транзакции
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *userRepository { // создаём репозиторий для работы с пользователями (конструктор)
	return &userRepository{
		queries: queries,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, email, passwordHash string) (User, error) { // отделям логику создания пользвателя в отдельный метод
	params := db.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	}
	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.Time,
	}, nil
}

func (r *userRepository) GetUserByEmail( // отделяем логику получения пользователя по email в отдельный метод
	ctx context.Context,
	email string,
) (User, error) {

	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}

		return User{}, fmt.Errorf("get user by email: %w", err)
	}

	return User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.Time,
	}, nil
}

// TODO: реализуй UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

type tokenRepository struct { // для работы с токенами в транзакции
	queries *db.Queries
}

func NewTokenRepository(queries *db.Queries) *tokenRepository { // создаём репозиторий для работы с токенами (конструктор)
	return &tokenRepository{
		queries: queries,
	}
}

func (r *tokenRepository) CreateRefreshToken( // отделяем логику создания токена в отдельный метод
	ctx context.Context,
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) (RefreshToken, error) {

	params := db.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		},
	}

	token, err := r.queries.CreateRefreshToken(ctx, params)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("create refresh token: %w", err)
	}

	var revokedAt *time.Time
	if token.RevokedAt.Valid {
		t := token.RevokedAt.Time
		revokedAt = &t
	}

	return RefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt.Time,
		RevokedAt: revokedAt,
		CreatedAt: token.CreatedAt.Time,
	}, nil
}

func (r *tokenRepository) GetActiveToken(
	ctx context.Context,
	tokenHash string,
) (RefreshToken, error) {

	token, err := r.queries.GetActiveToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RefreshToken{}, ErrTokenNotFound
		}

		return RefreshToken{}, fmt.Errorf("get active token: %w", err)
	}

	var revokedAt *time.Time
	if token.RevokedAt.Valid {
		t := token.RevokedAt.Time
		revokedAt = &t
	}

	return RefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt.Time,
		RevokedAt: revokedAt,
		CreatedAt: token.CreatedAt.Time,
	}, nil
}

func (r *tokenRepository) RevokeToken( // отделяем логику отзыва токена в отдельный метод
	ctx context.Context,
	tokenHash string,
) error {

	rowsAffected, err := r.queries.RevokeToken(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTokenNotFound
	}

	return nil
}

// TODO: реализуй TokenRepository
func RegisterWithToken(
	ctx context.Context,
	pool *pgxpool.Pool,
	email, pwdHash, tokenHash string,
	expiresAt time.Time,
) error {

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// sqlc работает через эту конкретную транзакцию
	qtx := db.New(tx)

	// репозитории тоже теперь работают через эту транзакцию
	userRepo := NewUserRepository(qtx)
	tokenRepo := NewTokenRepository(qtx)

	user, err := userRepo.CreateUser(ctx, email, pwdHash)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	token, err := tokenRepo.CreateRefreshToken(
		ctx,
		user.ID,
		tokenHash,
		expiresAt,
	)
	if err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	fmt.Printf("User created: %+v\n", user)
	fmt.Printf("Refresh token created: %+v\n", token)

	return nil
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://dev:dev@localhost:5432/task06db"
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Println("connect error:", err)
		return
	}
	defer pool.Close()

	// sqlc
	queries := db.New(pool)

	// наши репозитории
	userRepo := NewUserRepository(queries)
	tokenRepo := NewTokenRepository(queries)

	// Делаем уникальные данные, чтобы повторный go run . не падал на UNIQUE.
	n := time.Now().UnixNano()

	email := fmt.Sprintf("bob%d@example.com", n)
	passwordHash := "some-password-hash"
	tokenHash := fmt.Sprintf("token-%d", n)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	// Создаём пользователя + токен в одной транзакции.
	err = RegisterWithToken(
		ctx,
		pool,
		email,
		passwordHash,
		tokenHash,
		expiresAt,
	)
	if err != nil {
		fmt.Println("register error:", err)
		return
	}

	fmt.Println("register with token: ok")

	// Проверяем UserRepository.
	user, err := userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		fmt.Println("get user error:", err)
		return
	}

	fmt.Printf("User from repository: %+v\n", user)

	// Проверяем TokenRepository.
	token, err := tokenRepo.GetActiveToken(ctx, tokenHash)
	if err != nil {
		fmt.Println("get token error:", err)
		return
	}

	fmt.Printf("Active token: %+v\n", token)

	// Отзываем токен.
	err = tokenRepo.RevokeToken(ctx, tokenHash)
	if err != nil {
		fmt.Println("revoke token error:", err)
		return
	}

	fmt.Println("token revoked: ok")
}
