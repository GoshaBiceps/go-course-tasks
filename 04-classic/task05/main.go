// ============================================================
// Задача: Building H2O — LeetCode 1117  🔴 Senior
// ============================================================
//
// Задача с LeetCode. Частый вопрос на собесах уровня Senior.
//
// Есть два типа горутин: "водород" H и "кислород" O.
// Они должны формировать молекулы воды H2O строго:
//   - каждая молекула = 2 атома H + 1 атом O
//   - все три атома должны "встретиться" прежде чем какой-либо из них пройдёт дальше
//
// Реализуй:
//   type H2O struct { ... }
//   func NewH2O() *H2O
//   func (w *H2O) Hydrogen(fn func())   // fn() = "связаться в молекулу"
//   func (w *H2O) Oxygen(fn func())
//
// Ожидаемый вывод (буквы в группах по 3, каждая группа = HOO нет, = HHO или OHH):
//   OHH HHO OHH...  (по 2 H и 1 O в каждой тройке)
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

type H2O struct {
	hSem chan struct{} // пропускает до 2 водородов
	oSem chan struct{} // пропускает 1 кислород
	bar  *barrier
}

// barrier — встреча трёх горутин перед продолжением
type barrier struct {
	mu    sync.Mutex
	cond  *sync.Cond
	count int
	total int
}

func newBarrier(n int) *barrier {
	b := &barrier{total: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// TODO: реализуй Wait
// Подсказка: последний пришедший разбуждает всех, остальные ждут
func (b *barrier) Wait() {
	b.mu.Lock()
	// TODO
	b.mu.Unlock()
}

// TODO: реализуй NewH2O
// Подсказка: семафоры ограничивают сколько атомов каждого типа собирается в одну "встречу",
// а барьер синхронизирует их — подумай какие ёмкости нужны для H и O
func NewH2O() *H2O {
	return &H2O{}
}

// TODO: реализуй Hydrogen
func (w *H2O) Hydrogen(fn func()) {
	// TODO
}

// TODO: реализуй Oxygen
func (w *H2O) Oxygen(fn func()) {
	// TODO
}

// === Тесты ===

func TestH2O(t *testing.T) {
	const molecules = 10
	const total = molecules * 3 // 10 * (2H + 1O) = 30 атомов

	h2o := NewH2O()
	var mu sync.Mutex
	var result strings.Builder
	var wg sync.WaitGroup

	// 2*molecules водородов
	for range molecules * 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h2o.Hydrogen(func() {
				mu.Lock()
				result.WriteRune('H')
				mu.Unlock()
			})
		}()
	}

	// molecules кислородов
	for range molecules {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h2o.Oxygen(func() {
				mu.Lock()
				result.WriteRune('O')
				mu.Unlock()
			})
		}()
	}

	wg.Wait()

	s := result.String()
	if len(s) != total {
		t.Fatalf("ожидали %d символов, получили %d", total, len(s))
	}

	// Проверяем что в каждой тройке ровно 2 H и 1 O
	hCount := strings.Count(s, "H")
	oCount := strings.Count(s, "O")
	if hCount != molecules*2 {
		t.Errorf("H count = %d, want %d", hCount, molecules*2)
	}
	if oCount != molecules {
		t.Errorf("O count = %d, want %d", oCount, molecules)
	}
}

func main() {
	h2o := NewH2O()
	var mu sync.Mutex
	var result strings.Builder
	var wg sync.WaitGroup

	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h2o.Hydrogen(func() { mu.Lock(); result.WriteRune('H'); mu.Unlock() })
		}()
	}
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h2o.Oxygen(func() { mu.Lock(); result.WriteRune('O'); mu.Unlock() })
		}()
	}

	wg.Wait()
	fmt.Println(result.String()) // должно быть 2 молекулы H2O
}
