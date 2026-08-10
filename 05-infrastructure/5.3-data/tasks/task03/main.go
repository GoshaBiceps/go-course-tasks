// Задание: Репозиторий пользователей на pgx
//
// Реализуй UserRepository с методами Create и GetByID на реальном pgxpool.
// При отсутствии записи возвращай ErrUserNotFound.
//
// Для запуска нужен PostgreSQL. Пример docker-compose:
//   services:
//     postgres:
//       image: postgres:16
//       environment:
//         POSTGRES_USER: dev
//         POSTGRES_PASSWORD: dev
//         POSTGRES_DB: devdb
//       ports: ["5432:5432"]
//
// Таблица (создай вручную или миграцией из task02):
//   CREATE TABLE users (
//       id         bigserial primary key,
//       email      text not null unique,
//       created_at timestamptz not null default now()
//   );
//
// Ожидаемый результат:
//   created: {1 alice@example.com 2024-...}
//   found: {1 alice@example.com 2024-...}
//   not found: user not found

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID        int64
	Email     string
	CreatedAt time.Time
}

type UserRepository interface {
	Create(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
}

// TODO: реализуй pgxUserRepository, который хранит *pgxpool.Pool
type pgxUserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &pgxUserRepository{pool: pool}
}

func (r *pgxUserRepository) Create(ctx context.Context, email string) (User, error) {
	// TODO: выполни INSERT ... RETURNING id, email, created_at
	// Используй: row := r.pool.QueryRow(ctx, "INSERT INTO users...", email)
	// Затем row.Scan(&user.ID, &user.Email, &user.CreatedAt)
	var user User

	err := r.pool.QueryRow(
		ctx,
		"INSERT INTO users (email)  VALUES ($1) RETURNING id, email, created_at",
		email,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *pgxUserRepository) GetByID(ctx context.Context, id int64) (User, error) {
	// TODO: выполни SELECT id, email, created_at FROM users WHERE id = $1
	// Если pgx.ErrNoRows — верни ErrUserNotFound
	// Подсказка: errors.Is(err, pgx.ErrNoRows)
	var user User

	err := r.pool.QueryRow(
		ctx,
		"SELECT id, email, created_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://dev:dev@localhost:5432/devdb"
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Println("connect error:", err)
		return
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		fmt.Println("ping error:", err)
		return
	}

	repo := NewUserRepository(pool)

	u, err := repo.Create(ctx, "alice@example.com")
	if err != nil {
		fmt.Println("create error:", err)
		return
	}
	fmt.Printf("created: %+v\n", u)

	found, err := repo.GetByID(ctx, u.ID)
	if err != nil {
		fmt.Println("get error:", err)
	} else {
		fmt.Printf("found: %+v\n", found)
	}

	_, err = repo.GetByID(ctx, 99999)
	if errors.Is(err, ErrUserNotFound) {
		fmt.Println("not found:", err)
	}
}
