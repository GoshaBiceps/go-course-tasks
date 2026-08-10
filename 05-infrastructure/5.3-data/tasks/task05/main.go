// Задание: Транзакционная операция — Transfer
//
// Реализуй Transfer(ctx, fromID, toID, amount) в одной транзакции.
// При ошибке выполняй rollback; возвращай контекстные ошибки через %w.
//
// Ожидаемый результат:
//   transfer 300 from Alice to Bob: ok
//   Alice after: 700
//   Bob after:   800
//   transfer 2000 (too much): insufficient funds
//   transfer from unknown account: account not found

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type Account struct {
	ID      int64
	Owner   string
	Balance int64
}

// TODO: реализуй Transfer(ctx, pool, fromID, toID, amount int64) error
// Шаги:
//   1. Начать транзакцию: tx, err := pool.Begin(ctx)
//   2. Отложить rollback: defer tx.Rollback(ctx)
//   3. Получить from-аккаунт через SELECT ... FOR UPDATE (блокировка строки)
//   4. Проверить баланс >= amount, иначе вернуть ErrInsufficientFunds
//   5. Получить to-аккаунт
//   6. Обновить оба баланса через UPDATE
//   7. Выполнить tx.Commit(ctx)
// Подсказка: pgx.ErrNoRows → ErrAccountNotFound через errors.Is

func Transfer(
	ctx context.Context,
	pool *pgxpool.Pool,
	fromID, toID, amount int64,
) error {
	// Открываем транзакцию.
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Отложенный откат транзакции в случае ошибки.
	defer tx.Rollback(ctx)

	// Здесь будут храниться данные первого аккаунта.
	var fromAccount Account

	// Получаем from-аккаунт с блокировкой строки.
	err = tx.QueryRow(
		ctx,
		`SELECT id, owner, balance
		 FROM accounts
		 WHERE id = $1
		 FOR UPDATE`,
		fromID,
	).Scan(
		&fromAccount.ID,
		&fromAccount.Owner,
		&fromAccount.Balance,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Если аккаунт не найден, возвращаем ErrAccountNotFound.
			return fmt.Errorf(
				"select from account %d: %w",
				fromID,
				ErrAccountNotFound,
			)
		}

		// Если произошла другая ошибка, возвращаем её.
		return fmt.Errorf("select from account %d: %w", fromID, err)
	}

	// Проверяем, достаточно ли денег на балансе первого аккаунта.
	if fromAccount.Balance < amount {
		// Оборачиваем нашу ошибку ErrInsufficientFunds.
		return fmt.Errorf(
			"account %d has insufficient funds: %w",
			fromID,
			ErrInsufficientFunds,
		)
	}

	// Здесь будут храниться данные второго аккаунта.
	var toAccount Account

	// Получаем to-аккаунт с блокировкой строки.
	err = tx.QueryRow(
		ctx,
		`SELECT id, owner, balance
		 FROM accounts
		 WHERE id = $1
		 FOR UPDATE`,
		toID,
	).Scan(
		&toAccount.ID,
		&toAccount.Owner,
		&toAccount.Balance,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Если второй аккаунт не найден, возвращаем ErrAccountNotFound.
			return fmt.Errorf(
				"select to account %d: %w",
				toID,
				ErrAccountNotFound,
			)
		}

		// Если произошла другая ошибка, возвращаем её.
		return fmt.Errorf("select to account %d: %w", toID, err)
	}

	// Обновляем данные в базе: уменьшаем баланс первого аккаунта.
	_, err = tx.Exec(
		ctx,
		`UPDATE accounts
		 SET balance = balance - $1
		 WHERE id = $2`,
		amount,
		fromID,
	)

	if err != nil {
		return fmt.Errorf(
			"decrease balance for account %d: %w",
			fromID,
			err,
		)
	}

	// Обновляем данные в базе: увеличиваем баланс второго аккаунта.
	_, err = tx.Exec(
		ctx,
		`UPDATE accounts
		 SET balance = balance + $1
		 WHERE id = $2`,
		amount,
		toID,
	)

	if err != nil {
		return fmt.Errorf(
			"increase balance for account %d: %w",
			toID,
			err,
		)
	}

	// Сохраняем результаты транзакции.
	// После успешного Commit изменения станут видны другим операциям.
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transfer: %w", err)
	}

	// Ошибок нет — перевод успешно выполнен.
	return nil
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

	// Предполагается, что в БД уже есть таблица accounts и тестовые данные
	// CREATE TABLE accounts (id bigserial primary key, owner text, balance bigint);
	// INSERT INTO accounts (owner, balance) VALUES ('Alice', 1000), ('Bob', 500);

	if err = Transfer(ctx, pool, 1, 2, 300); err != nil {
		fmt.Println("transfer 300 from Alice to Bob:", err)
	} else {
		fmt.Println("transfer 300 from Alice to Bob: ok")
	}

	if err = Transfer(ctx, pool, 1, 2, 2000); err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			fmt.Println("transfer 2000 (too much): insufficient funds")
		} else {
			fmt.Println("unexpected error:", err)
		}
	}

	if err = Transfer(ctx, pool, 9999, 2, 100); err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			fmt.Println("transfer from unknown account: account not found")
		} else {
			fmt.Println("unexpected error:", err)
		}
	}
}
