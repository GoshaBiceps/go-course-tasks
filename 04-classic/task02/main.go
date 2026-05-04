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
	"testing"
)

// === Вариант 1: через каналы ===

// TODO: реализуй producerConsumerChan
// Подсказка: два буферизованных канала и два WaitGroup — для производителей и потребителей
func producerConsumerChan(producers, consumers, n, bufSize int) []int {
	return nil
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
	return nil
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
