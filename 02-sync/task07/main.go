// ============================================================
// Задача: RWMutex с приоритетом писателей  🔴 Senior
// ============================================================
//
// Стандартный sync.RWMutex не гарантирует что писатель когда-нибудь
// получит доступ при постоянном потоке читателей — возможна starvation
// (хотя в Go на практике голодания нет, это контракт не гарантирует).
//
// Реализуй СВОЙ RWMutex с явной гарантией: если писатель хочет получить
// блокировку, новые читатели не могут войти в критическую секцию пока
// писатель её не отработает.
//
//   type WriterPriorityRWMutex struct { ... }
//
//   func New() *WriterPriorityRWMutex
//   func (m *WriterPriorityRWMutex) RLock()
//   func (m *WriterPriorityRWMutex) RUnlock()
//   func (m *WriterPriorityRWMutex) Lock()
//   func (m *WriterPriorityRWMutex) Unlock()
//
// Алгоритм (примерный):
//   - счётчик активных читателей
//   - счётчик ожидающих писателей
//   - RLock: если есть ожидающие писатели — жди; иначе увеличь reader count
//   - Lock: увеличь waiting writers, жди пока reader count == 0 и нет активного писателя
//
// Используй sync.Mutex + sync.Cond.
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"fmt"
	"sync"
	"time"
)

type WriterPriorityRWMutex struct {
	mu             sync.Mutex
	readerCount    int
	writerActive   bool
	writersWaiting int
	readerCond     *sync.Cond
	writerCond     *sync.Cond
}

// TODO: реализуй конструктор
func New() *WriterPriorityRWMutex {
	m := &WriterPriorityRWMutex{}
	m.readerCond = sync.NewCond(&m.mu)
	m.writerCond = sync.NewCond(&m.mu)
	return m
}

// TODO: реализуй RLock
// Подсказка: читатель НЕ должен войти если есть активный писатель или если
// есть ожидающие писатели (это и есть "приоритет писателей")
func (m *WriterPriorityRWMutex) RLock() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// TODO
}

// TODO: реализуй RUnlock
// Подсказка: когда последний читатель уходит — надо разбудить ожидающего писателя
func (m *WriterPriorityRWMutex) RUnlock() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// TODO
}

// TODO: реализуй Lock
func (m *WriterPriorityRWMutex) Lock() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// TODO: зарегистрируй ожидание, жди пока читатели уйдут и нет другого писателя
}

// TODO: реализуй Unlock
// Подсказка: когда писатель выходит — кого будить? И тех и других?
// Подумай какие семантики Signal/Broadcast подходят
func (m *WriterPriorityRWMutex) Unlock() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// TODO
}

func main() {
	rw := New()
	var wg sync.WaitGroup

	// Постоянный поток читателей
	for i := 0; i < 5; i++ {
		wg.Add(1)
		id := i
		go func() {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				rw.RLock()
				fmt.Printf("читатель %d — заход %d\n", id, j)
				time.Sleep(20 * time.Millisecond)
				rw.RUnlock()
				time.Sleep(5 * time.Millisecond)
			}
		}()
	}

	// Писатель приходит посреди потока читателей
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(15 * time.Millisecond)
		rw.Lock()
		fmt.Println(">>> писатель — работает")
		time.Sleep(50 * time.Millisecond)
		fmt.Println(">>> писатель — готово")
		rw.Unlock()
	}()

	wg.Wait()
}
