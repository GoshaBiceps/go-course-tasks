// Задание 1: Безопасный счётчик через Mutex
//
// Создай структуру SafeCounter с:
//   - полем mu типа sync.Mutex
//   - полем value типа int
//   - методом Increment() - увеличивает value на 1 (с блокировкой)
//   - методом Value() int - возвращает value (с блокировкой)
//
// Запусти 1000 горутин, каждая вызывает counter.Increment().
// После wg.Wait() выведи финальное значение.
//
// Проверь отсутствие гонок:
//   go run -race main.go
//
// Ожидаемый вывод:
//   Финальный счётчик: 1000
//   (и никаких WARNING: DATA RACE)
//
// Запусти: go run main.go

package main

import (
	"fmt"
	"sync"
)

// TODO: напиши структуру SafeCounter
// type SafeCounter struct { ... }
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

// TODO: напиши метод Increment()
func (s *SafeCounter) Increment() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.value++

}

// TODO: напиши метод Value() int
func (s *SafeCounter) Value() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.value
}

func main() {
	// TODO: создай SafeCounter и запусти 1000 горутин
	wg := sync.WaitGroup{}
	counter := SafeCounter{
		mu:    sync.Mutex{},
		value: 0,
	}

	// Каждая горутина вызывает counter.Increment()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	// После завершения всех горутин выведи counter.Value()
	fmt.Println("Финальный счётчик:", counter.value)

}
