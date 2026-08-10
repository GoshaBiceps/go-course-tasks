package main

import (
	"context"
	"errors"
	"fmt"
)

type Account struct {
	ID      int64
	Owner   string
	Balance int64
}

var ErrAccountNotFound = errors.New("account not found")

type AccountRepository interface {
	Create(ctx context.Context, account Account) (Account, error)
	GetByID(ctx context.Context, id int64) (Account, error)
}

// TODO: Implement InMemoryAccountRepository.
// Use map[int64]Account as storage and an auto-incrementing ID counter.
// GetByID must return ErrAccountNotFound when the key is missing.

type InMemoryAccountRepository struct {
	accounts map[int64]Account
	nextID   int64
}

func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		accounts: make(map[int64]Account),
		nextID:   1,
	}
}

func (r *InMemoryAccountRepository) Create(ctx context.Context, account Account) (Account, error) {
	// TODO: assign r.nextID to account.ID, store in map, increment counter, return account

	account.ID = r.nextID            // так добавляем уникальный ID к аккаунту
	r.accounts[account.ID] = account // сохраняем аккаунт в мапу
	r.nextID++                       // увеличиваем счетчик ID для следующего аккаунта

	return account, nil
}

func (r *InMemoryAccountRepository) GetByID(ctx context.Context, id int64) (Account, error) {
	// TODO: look up id in r.accounts; return ErrAccountNotFound if missing
	if account, exsists := r.accounts[id]; exsists { // проверяем, есть ли аккаунт с таким ID в мапе
		return account, nil // если есть, возвращаем его
	}
	return Account{}, ErrAccountNotFound // если нет , возвращаем ошибку
}

func main() {
	repo := NewInMemoryAccountRepository()
	ctx := context.Background()

	a1, err := repo.Create(ctx, Account{Owner: "Alice", Balance: 1000})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Created: %+v\n", a1)

	found, err := repo.GetByID(ctx, a1.ID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Found: %+v\n", found)
	}

	_, err = repo.GetByID(ctx, 999)
	if errors.Is(err, ErrAccountNotFound) {
		fmt.Println("Account 999: not found (expected)")
	}
}
