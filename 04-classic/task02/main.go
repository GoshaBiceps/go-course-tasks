// ============================================================
// Задача: Producer-Consumer с bounded buffer  🟡 Middle
// ============================================================
//
// Классика на собесах Junior/Middle уровня.
//
// Реализуй через каналы:
//   - M производителей генерируют числа 0..N
//   - K потребителей читают, возводят в квадрат, пишут в results
//   - Буфер между ними ограничен (размер B)
//
// Требования:
//   - Потребители завершаются когда производители закончили И буфер пуст
//   - Нет утечек горутин
//   - Все числа должны быть обработаны ровно один раз
//
// Реализуй ДВА варианта:
//   1. Через каналы (идиоматично в Go)
//   2. Через sync.Cond (для понимания классических примитивов)
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"fmt"
	"sort"
	"sync"
	"testing"
)

// === Вариант 1: через каналы ===

// TODO: реализуй producerConsumerChan
// Подсказка: два буферизованных канала и два WaitGroup — для производителей и потребителей
func producerConsumerChan(producers, consumers, n, bufSize int) []int {
	jobs := make(chan int, bufSize)
	results := make(chan int, n)

	var producerWG sync.WaitGroup
	var consumerWG sync.WaitGroup

	// producers
	for p := 0; p < producers; p++ {
		producerWG.Add(1)

		go func(start int) {
			defer producerWG.Done()

			for i := start; i < n; i += producers {
				jobs <- i
			}
		}(p)
	}

	// consumers
	for c := 0; c < consumers; c++ {
		consumerWG.Add(1)

		go func() {
			defer consumerWG.Done()

			for v := range jobs {
				results <- v * v
			}
		}()
	}

	// закрываем jobs после producers
	go func() {
		producerWG.Wait()
		close(jobs)
	}()

	// закрываем results после consumers
	go func() {
		consumerWG.Wait()
		close(results)
	}()

	// собираем результаты
	var out []int

	for v := range results {
		out = append(out, v)
	}

	return out
}

func TestProducerConsumer(t *testing.T) {
	results := producerConsumerChan(3, 4, 20, 5)
	sort.Ints(results)

	if len(results) != 20 {
		t.Fatalf("ожидали 20 результатов, получили %d", len(results))
	}

	// Проверяем что это квадраты чисел 0..19
	for i, v := range results {
		want := i * i
		if v != want {
			t.Errorf("[%d] = %d, want %d", i, v, want)
		}
	}
}

// === Вариант 2: через sync.Cond ===

// TODO: реализуй producerConsumerCond
// Подсказка: буфер — обычный срез; производители ждут пока буфер полон, потребители — пока пуст
func producerConsumerCond(producers, consumers, n, bufSize int) []int {
	var (
		mu       sync.Mutex
		notEmpty = sync.NewCond(&mu)
		notFull  = sync.NewCond(&mu)

		buffer  []int
		results []int

		producerWG sync.WaitGroup
		consumerWG sync.WaitGroup

		producersDone int
	)

	// producers
	for p := 0; p < producers; p++ {
		producerWG.Add(1)

		go func(start int) {
			defer producerWG.Done()

			for i := start; i < n; i += producers {
				mu.Lock()

				// ждём пока появится место
				for len(buffer) == bufSize {
					notFull.Wait()
				}

				buffer = append(buffer, i)

				// будим consumer
				notEmpty.Signal()

				mu.Unlock()
			}
		}(p)
	}

	// когда producers закончились
	go func() {
		producerWG.Wait()

		mu.Lock()
		producersDone = producers

		// будим всех consumers
		notEmpty.Broadcast()
		mu.Unlock()
	}()

	// consumers
	for c := 0; c < consumers; c++ {
		consumerWG.Add(1)

		go func() {
			defer consumerWG.Done()

			for {
				mu.Lock()

				// ждём пока buffer пуст
				for len(buffer) == 0 {
					// если producers закончились
					// и buffer пуст -> выходим
					if producersDone == producers {
						mu.Unlock()
						return
					}

					notEmpty.Wait()
				}

				v := buffer[0]
				buffer = buffer[1:]

				// освободилось место
				notFull.Signal()

				mu.Unlock()

				sq := v * v

				mu.Lock()
				results = append(results, sq)
				mu.Unlock()
			}
		}()
	}

	consumerWG.Wait()

	return results
}

func TestProducerConsumerCond(t *testing.T) {
	results := producerConsumerCond(3, 4, 20, 5)
	sort.Ints(results)

	if len(results) != 20 {
		t.Fatalf("ожидали 20 результатов, получили %d", len(results))
	}
	for i, v := range results {
		if v != i*i {
			t.Errorf("[%d] = %d, want %d", i, v, i*i)
		}
	}
}

func main() {
	results := producerConsumerChan(2, 3, 10, 3)
	sort.Ints(results)
	fmt.Println("Результаты:", results)
}
