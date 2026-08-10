// Задание: Слой данных через sqlc
//
// Создай SQL-запросы и сгенерируй type-safe код через sqlc.
//
// Структура файлов (после генерации):
//   task04/
//     sqlc.yaml
//     sql/
//       schema.sql
//       queries.sql
//     db/           ← сгенерированный код
//       db.go
//       models.go
//       queries.sql.go
//     main.go
//
// Шаги:
// 1. Установи sqlc: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
// 2. Создай sqlc.yaml (см. комментарий ниже)
// 3. Создай sql/schema.sql и sql/queries.sql
// 4. Запусти: sqlc generate
// 5. Реализуй main, используя сгенерированный код
//
// sqlc.yaml (пример):
//   version: "2"
//   sql:
//     - engine: "postgresql"
//       queries: "sql/queries.sql"
//       schema:  "sql/schema.sql"
//       gen:
//         go:
//           package: "db"
//           out: "db"
//
// sql/queries.sql (TODO: создай файл с этими запросами):
//
//   -- name: CreateUser :one
//   INSERT INTO users (email) VALUES ($1) RETURNING *;
//
//   -- name: GetUser :one
//   SELECT * FROM users WHERE id = $1;
//
//   -- name: ListUsers :many
//   SELECT * FROM users ORDER BY id;
//
// Ожидаемый результат после реализации:
//   created: {ID:1 Email:bob@example.com CreatedAt:...}
//   found: {ID:1 Email:bob@example.com CreatedAt:...}
//   all users: [{...}]

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cource/data-task04/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

	queries := db.New(pool)

	u, err := queries.CreateUser(ctx, "bob@example.com")
	if err != nil {
		fmt.Println("create error:", err)
		return
	}
	fmt.Printf("created: %+v\n", u)

	found, err := queries.GetUser(ctx, u.ID)
	if err != nil {
		fmt.Println("get error:", err)
	} else {
		fmt.Printf("found: %+v\n", found)
	}

	users, err := queries.ListUsers(ctx)
	if err != nil {
		fmt.Println("list error:", err)
	} else {
		fmt.Println("all users:", users)
	}

}
