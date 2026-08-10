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

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type AccountRepository interface {
	Create(ctx context.Context, account Account) (Account, error)
	GetByID(ctx context.Context, id int64) (Account, error)
	Update(ctx context.Context, account Account) error
}

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
	account.ID = r.nextID
	r.accounts[r.nextID] = account
	r.nextID++
	return account, nil
}

func (r *InMemoryAccountRepository) GetByID(ctx context.Context, id int64) (Account, error) {
	a, ok := r.accounts[id]
	if !ok {
		return Account{}, ErrAccountNotFound
	}
	return a, nil
}

func (r *InMemoryAccountRepository) Update(ctx context.Context, account Account) error {
	if _, ok := r.accounts[account.ID]; !ok {
		return ErrAccountNotFound
	}
	r.accounts[account.ID] = account
	return nil
}

// TODO: Implement AccountService.Transfer(ctx, fromID, toID, amount int64) error
// 1. Get both accounts via s.repo.GetByID (propagate ErrAccountNotFound)
// 2. Check from.Balance >= amount (return ErrInsufficientFunds otherwise)
// 3. Debit from, credit to
// 4. Update both accounts via s.repo.Update

type AccountService struct {
	repo AccountRepository
}

func NewAccountService(repo AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Transfer(ctx context.Context, fromID, toID, amount int64) error {
	// TODO: implement
	accountFrom, err := s.repo.GetByID(ctx, fromID) // получаем аккаунт отправителя
	if err != nil {
		return err // если аккаунт не найден, возвращаем ошибку
	}

	accountTo, err := s.repo.GetByID(ctx, toID) // получаем аккаунт получателя
	if err != nil {
		return err // если аккаунт получателя не найден, возвращаем ошибку
	}

	if accountFrom.Balance < amount { // проверяем, достаточно ли средств на аккаунте отправителя
		return ErrInsufficientFunds // если средств недостаточно, возвращаем ошибку
	}

	accountFrom.Balance -= amount // списываем средства с аккаунта отправителя
	accountTo.Balance += amount   // зачисляем средства на аккаунт получателя

	if err := s.repo.Update(ctx, accountFrom); err != nil { // обновляем аккаунт отправителя
		return err // если ошибка при обновлении, возвращаем её
	}
	if err := s.repo.Update(ctx, accountTo); err != nil { // обновляем аккаунт получателя
		return err // если ошибка при обновлении, возвращаем её
	}

	return nil // если всё прошло успешно, возвращаем nil
}

func main() {
	repo := NewInMemoryAccountRepository()
	ctx := context.Background()

	a1, _ := repo.Create(ctx, Account{Owner: "Alice", Balance: 1000})
	a2, _ := repo.Create(ctx, Account{Owner: "Bob", Balance: 500})

	svc := NewAccountService(repo)

	fmt.Println("=== Transfer 300 from Alice to Bob ===")
	if err := svc.Transfer(ctx, a1.ID, a2.ID, 300); err != nil {
		fmt.Println("Error:", err)
	}
	alice, _ := repo.GetByID(ctx, a1.ID)
	bob, _ := repo.GetByID(ctx, a2.ID)
	fmt.Printf("Alice: %+v\n", alice)
	fmt.Printf("Bob: %+v\n", bob)

	fmt.Println("\n=== Transfer 2000 from Alice to Bob (should fail: insufficient funds) ===")
	if err := svc.Transfer(ctx, a1.ID, a2.ID, 2000); err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\n=== Transfer from non-existent account ===")
	if err := svc.Transfer(ctx, 999, a2.ID, 100); err != nil {
		fmt.Println("Error:", err)
	}
}
