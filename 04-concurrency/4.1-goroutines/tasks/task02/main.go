// Задание 2: Найди и исправь гонку данных
//
// В этом коде несколько горутин одновременно увеличивают переменную counter.
// Это называется "гонка данных" - результат непредсказуем.
//
// Шаг 1: Запусти с флагом -race и прочитай сообщение об ошибке:
//   go run -race main.go
//
// Ты увидишь что-то вроде:
//   WARNING: DATA RACE
//   Write at 0x... by goroutine ...
//   Previous write at 0x... by goroutine ...
//
// Шаг 2: Исправь код, используя sync/atomic для атомарного увеличения:
//   import "sync/atomic"
//   var counter int64
//   atomic.AddInt64(&counter, 1)  // вместо counter++
//
// После исправления: go run -race main.go не должен выдавать WARNING.
//
// Запусти: go run main.go (или go run -race main.go)

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter atomic.Int64 // ПРОБЛЕМА: несколько горутин пишут сюда одновременно
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			counter.Add(1)
			defer wg.Done()

		}()
	}

	wg.Wait()
	// Из-за гонки counter может быть меньше 1000
	fmt.Println("Финальное значение счётчика:", counter)
	fmt.Println("(Правильное значение должно быть 1000)")
}
