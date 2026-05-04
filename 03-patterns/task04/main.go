// ============================================================
// Задача: Future / Promise  🔴 Senior
// ============================================================
//
// Реализуй паттерн Future — асинхронное вычисление результата.
//
//   func Async[T any](fn func() (T, error)) *Future[T]
//
//   type Future[T any] struct { ... }
//   func (f *Future[T]) Await() (T, error)                    // блокирует до результата
//   func (f *Future[T]) AwaitTimeout(d time.Duration) (T, error, bool) // bool = ok
//   func (f *Future[T]) Then(fn func(T) T) *Future[T]         // цепочка трансформаций
//   func (f *Future[T]) Done() <-chan struct{}                  // закрывается когда готово
//
// Требования:
//   - fn запускается сразу в отдельной горутине
//   - Await можно вызывать из нескольких горутин одновременно
//   - Then возвращает новый Future (результат доступен после обоих вычислений)
//
// Пример:
//   f := Async(func() (int, error) {
//       time.Sleep(100 * time.Millisecond)
//       return 42, nil
//   }).Then(func(v int) int { return v * 2 })
//
//   result, err := f.Await() // result = 84

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrFutureTimeout = errors.New("future: таймаут истёк")

type Future[T any] struct {
	once  sync.Once
	done  chan struct{}
	value T
	err   error
}

// TODO: реализуй Async — запускает fn в горутине, возвращает Future
// Подсказка: закрытие канала — универсальный сигнал готовности
func Async[T any](fn func() (T, error)) *Future[T] {
	return nil
}

// TODO: реализуй Await
func (f *Future[T]) Await() (T, error) {
	var zero T
	return zero, nil
}

// TODO: реализуй AwaitTimeout
func (f *Future[T]) AwaitTimeout(d time.Duration) (T, error, bool) {
	var zero T
	return zero, ErrFutureTimeout, false
}

// TODO: реализуй Done
func (f *Future[T]) Done() <-chan struct{} {
	return nil
}

// TODO: реализуй Then — цепочка: дождись текущего Future и примени fn к результату
func (f *Future[T]) Then(fn func(T) T) *Future[T] {
	return nil
}

func main() {
	// Простой Future
	f1 := Async(func() (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	// Цепочка
	f2 := f1.Then(func(v int) int { return v * 2 })
	f3 := f2.Then(func(v int) int { return v + 8 })

	result, err := f3.Await()
	fmt.Printf("42 * 2 + 8 = %d, err = %v\n", result, err) // 92

	// Таймаут
	slow := Async(func() (string, error) {
		time.Sleep(1 * time.Second)
		return "done", nil
	})

	_, _, ok := slow.AwaitTimeout(50 * time.Millisecond)
	fmt.Printf("таймаут сработал: %v\n", !ok) // true

	// Параллельное ожидание
	var wg sync.WaitGroup
	f4 := Async(func() (int, error) {
		time.Sleep(50 * time.Millisecond)
		return 100, nil
	})
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			v, _ := f4.Await()
			fmt.Printf("горутина %d получила: %d\n", id, v)
		}(i)
	}
	wg.Wait()
}
