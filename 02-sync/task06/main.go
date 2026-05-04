// ============================================================
// Задача: TryLock Mutex  🟡 Middle
// ============================================================
//
// Реализуй свой мьютекс на каналах (не используй sync.Mutex).
//
//   type TryMutex struct { ... }
//
//   func NewTryMutex() *TryMutex
//   func (m *TryMutex) Lock()
//   func (m *TryMutex) Unlock()
//   func (m *TryMutex) TryLock() bool
//   func (m *TryMutex) LockTimeout(d time.Duration) bool
//   func (m *TryMutex) LockContext(ctx context.Context) error
//
// Требования:
//   - Lock блокирует пока не захватит
//   - TryLock — non-blocking: true если захватил, false если уже занят
//   - LockTimeout — ждёт максимум d; false если не успел
//   - LockContext — отменяется через ctx
//   - Unlock без удерживания — паника с "unlock of unlocked mutex"
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TryMutex struct {
	ch chan struct{}
}

// TODO: реализуй NewTryMutex
// Подсказка: сам факт "владения" можно выразить наличием токена в канале
func NewTryMutex() *TryMutex {
	return nil
}

// TODO: реализуй Lock
func (m *TryMutex) Lock() {
}

// TODO: реализуй TryLock — non-blocking
// Подсказка: как через select понять что канал "не готов прямо сейчас"?
func (m *TryMutex) TryLock() bool {
	return false
}

// TODO: реализуй LockTimeout
func (m *TryMutex) LockTimeout(d time.Duration) bool {
	return false
}

// TODO: реализуй LockContext
func (m *TryMutex) LockContext(ctx context.Context) error {
	return nil
}

// TODO: реализуй Unlock (с паникой при двойном Unlock)
func (m *TryMutex) Unlock() {
}

func main() {
	m := NewTryMutex()

	m.Lock()
	fmt.Println("захватил")

	if !m.TryLock() {
		fmt.Println("TryLock вернул false — уже занят")
	}

	ok := m.LockTimeout(50 * time.Millisecond)
	fmt.Println("LockTimeout после занятого:", ok) // false

	m.Unlock()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		id := i
		go func() {
			defer wg.Done()
			m.Lock()
			fmt.Printf("горутина %d в критической секции\n", id)
			time.Sleep(20 * time.Millisecond)
			m.Unlock()
		}()
	}
	wg.Wait()

	// Проверка контекста
	m.Lock()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	err := m.LockContext(ctx)
	fmt.Println("LockContext err:", err) // context deadline exceeded
	m.Unlock()
}
