// ============================================================
// Задача: Модель Акторов  ⚫ Expert
// ============================================================
//
// Вопрос с финальных этапов собеса уровня Staff+.
//
// Модель акторов — альтернатива мьютексам.
// Актор:
//   - Имеет приватное состояние (никакого прямого доступа снаружи)
//   - Общается только через сообщения
//   - Обрабатывает сообщения строго последовательно (нет гонок)
//
// Реализуй:
//
//   type Actor[S any] struct { ... }
//
//   func NewActor[S any](initial S) *Actor[S]
//   func (a *Actor[S]) Send(msg func(state *S))  // отправить сообщение (не блокируется)
//   func (a *Actor[S]) Ask(msg func(state *S) any) any  // запрос с ответом (блокируется)
//   func (a *Actor[S]) Stop()
//
// Пример — счётчик без мьютекса:
//   counter := NewActor(0)
//   counter.Send(func(n *int) { *n++ })
//   result := counter.Ask(func(n *int) any { return *n })
//   fmt.Println(result) // 1
//
// Проверь что нет гонок:
//   go test -race -v ./...

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type message[S any] struct {
	fn      func(state *S)
	replyCh chan any
	ask     func(state *S) any
}

type Actor[S any] struct {
	state   S
	mailbox chan message[S]
	done    chan struct{}
	once    sync.Once
}

// TODO: реализуй NewActor
// Подсказка: запусти run() в горутине — она обрабатывает все сообщения последовательно
func NewActor[S any](initial S) *Actor[S] {
	return nil
}

// TODO: реализуй run — бесконечный цикл обработки сообщений из mailbox
// Подсказка: message может быть двух видов (fn или ask) — для ask нужно отправить результат в replyCh
// Цикл должен завершаться при закрытии a.done
func (a *Actor[S]) run() {
	// TODO
}

// TODO: реализуй Send — fire-and-forget; учитывай что actor может быть остановлен
func (a *Actor[S]) Send(fn func(state *S)) {
}

// TODO: реализуй Ask — запрос с ожиданием ответа
func (a *Actor[S]) Ask(fn func(state *S) any) any {
	return nil
}

// TODO: реализуй Stop
func (a *Actor[S]) Stop() {
}

// === Пример: Банковский счёт без мьютекса ===

type BankAccount struct {
	balance float64
	txCount int
}

func NewBankAccount(initial float64) *Actor[BankAccount] {
	return NewActor(BankAccount{balance: initial})
}

func Deposit(account *Actor[BankAccount], amount float64) {
	account.Send(func(s *BankAccount) {
		s.balance += amount
		s.txCount++
	})
}

func Withdraw(account *Actor[BankAccount], amount float64) bool {
	result := account.Ask(func(s *BankAccount) any {
		if s.balance < amount {
			return false
		}
		s.balance -= amount
		s.txCount++
		return true
	})
	return result.(bool)
}

func Balance(account *Actor[BankAccount]) float64 {
	result := account.Ask(func(s *BankAccount) any {
		return s.balance
	})
	return result.(float64)
}

func main() {
	account := NewBankAccount(1000.0)

	var wg sync.WaitGroup
	var successWithdrawals atomic.Int32

	// 100 горутин пытаются снять по 20
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if Withdraw(account, 20.0) {
				successWithdrawals.Add(1)
			}
		}()
	}

	// 50 горутин пополняют по 10
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Deposit(account, 10.0)
		}()
	}

	wg.Wait()
	account.Stop()

	fmt.Printf("Баланс: %.2f\n", Balance(account))
	fmt.Printf("Успешных снятий: %d\n", successWithdrawals.Load())
	// Баланс не должен уйти в минус!
}
