// ============================================================
// Задача: Cigarette Smokers Problem  🔴 Senior
// ============================================================
//
// Классическая задача Паттерсона (1977). Показывает ограничения
// "наивных" семафоров.
//
// За столом — 3 курильщика и 1 агент.
// Чтобы скрутить и выкурить сигарету нужно 3 ингредиента: табак, бумага, спички.
// У каждого курильщика бесконечный запас одного ингредиента:
//   Курильщик 1: табак
//   Курильщик 2: бумага
//   Курильщик 3: спички
// Агент кладёт на стол ДВА случайных ингредиента (из трёх возможных пар).
// Курильщик с НЕДОСТАЮЩИМ третьим — забирает оба, курит, сигнализирует агенту.
// Агент снова кладёт пару ингредиентов. И так далее.
//
// Реализуй (без deadlock, без busy-waiting):
//
//   type Table struct { ... }
//
//   func NewTable() *Table
//   func (t *Table) AgentRound()          // кладёт случайную пару ингредиентов
//   func (t *Table) Smoker(has Ingredient, onSmoke func()) // работа курильщика
//
// Требования:
//   - Только тот курильщик у кого нет этих ингредиентов — забирает пару
//   - Курильщики не могут общаться напрямую, только через стол
//   - Нет голодания — все три курильщика должны курить +- поровну
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Ingredient int

const (
	Tobacco Ingredient = iota
	Paper
	Matches
)

type Table struct {
	mu sync.Mutex
	// TODO: подумай какие поля нужны чтобы закодировать "что лежит на столе"
	// и как будить правильного курильщика
}

// TODO: реализуй конструктор
func NewTable() *Table {
	return &Table{}
}

// TODO: реализуй AgentRound
// Подсказка: выбери 2 случайных разных ингредиента, положи на стол,
// разбуди нужного курильщика; дождись пока он заберёт и вернёт стол в пустое состояние
func (t *Table) AgentRound() {
	// TODO
}

// TODO: реализуй Smoker
// Курильщик в бесконечном цикле:
//   1) ждёт пока на столе окажутся ДВА ингредиента, ни один из которых не его
//   2) забирает их и курит (вызывает onSmoke)
//   3) освобождает стол (сигнализирует агенту)
// Подсказка: по ингредиенту который есть у курильщика можно вычислить ПАРУ
// которую он ждёт. Используй sync.Cond или отдельные каналы на каждую пару.
func (t *Table) Smoker(has Ingredient, onSmoke func()) {
	// TODO
}

func main() {
	t := NewTable()
	var counts [3]atomic.Int32

	var wg sync.WaitGroup
	for i, ing := range []Ingredient{Tobacco, Paper, Matches} {
		wg.Add(1)
		idx := i
		has := ing
		go func() {
			defer wg.Done()
			for range 10 {
				t.Smoker(has, func() {
					counts[idx].Add(1)
					fmt.Printf("курильщик с %v курит\n", has)
					time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
				})
			}
		}()
	}

	go func() {
		for i := 0; i < 30; i++ {
			t.AgentRound()
		}
	}()

	wg.Wait()
	fmt.Printf("Итого: tobacco=%d paper=%d matches=%d\n",
		counts[0].Load(), counts[1].Load(), counts[2].Load())
}
