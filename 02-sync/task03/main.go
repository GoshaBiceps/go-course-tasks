// ============================================================
// Задача: Взвешенный семафор  🟡 Middle
// ============================================================
//
// Вопрос с собесов уровня Middle.
//
// Реализуй Semaphore с поддержкой "веса" (weighted semaphore):
//   - Ресурс имеет ёмкость N
//   - Acquire(n) захватывает n единиц. Блокируется если доступно < n.
//   - Release(n) возвращает n единиц.
//   - TryAcquire(n) — non-blocking: захватывает или возвращает false
//
// Примеры использования:
//   - Ограничение числа параллельных HTTP-запросов (каждый = 1 единица)
//   - Управление памятью (запрос на 10МБ = 10 единиц)
//   - Rate limiting по "стоимости" операции
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

type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore создаёт семафор с ёмкостью n.
// Подсказка: сам канал представляет токены
func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, n),
	}
}

// Acquire блокирующий захват n единиц.
// TODO: реализуй
func (s *Semaphore) Acquire(n int) {
	for i := 0; i < n; i++ {
		s.ch <- struct{}{}
	}
}

// AcquireContext захват с контекстом — можно отменить.
// TODO: если отмена настигнет в середине — верни уже захваченное и верни ошибку
func (s *Semaphore) AcquireContext(ctx context.Context, n int) error {
	acquired := 0

	for i := 0; i < n; i++ {
		select {
		case s.ch <- struct{}{}:
			acquired++
		case <-ctx.Done():

			for j := 0; j < acquired; j++ {
				<-s.ch
			}

			return ctx.Err()
		}
	}
	return nil
}

// TryAcquire non-blocking захват. Возвращает false если доступно < n.
// TODO: необходимо попытаться захватить не блокируясь; если не удалось — откатить всё что уже взяли
func (s *Semaphore) TryAcquire(n int) bool {
	acquired := 0

	for i := 0; i < n; i++ {
		select {
		case s.ch <- struct{}{}:
			acquired++

		default:
			for j := 0; j < acquired; j++ {
				<-s.ch
			}

			return false
		}
	}

	return true
}

// Release возвращает n единиц.
// TODO: реализуй
func (s *Semaphore) Release(n int) {
	for i := 0; i < n; i++ {
		<-s.ch
	}
}

// Available возвращает количество свободных единиц.
func (s *Semaphore) Available() int {
	return len(s.ch)
}

func main() {
	sem := NewSemaphore(3)
	var wg sync.WaitGroup

	// 5 задач, каждая занимает 1 единицу, не более 3 одновременно
	for i := 0; i < 5; i++ {
		wg.Add(1)
		n := i
		go func() {
			defer wg.Done()
			sem.Acquire(1)
			defer sem.Release(1)
			fmt.Printf("задача %d выполняется (доступно: %d)\n", n, sem.Available())
			time.Sleep(100 * time.Millisecond)
		}()
	}
	wg.Wait()
}
