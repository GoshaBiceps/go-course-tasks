// ============================================================
// Задача: Обедающие философы без дедлока  🔴 Senior
// ============================================================
//
// Классическая задача Дейкстры. Частый вопрос на собесах уровня Senior.
//
// Условие:
//   5 философов сидят за круглым столом.
//   Между каждыми двумя соседями — одна вилка (5 вилок всего).
//   Чтобы есть, философу нужны ДВЕ вилки: левая и правая.
//   Философ: думает → берёт вилки → ест → кладёт вилки → думает...
//
// Проблема: если все возьмут левую вилку одновременно — дедлок.
//
// Реализуй без дедлока. Три способа:
//
//   A) Resource hierarchy: нечётный философ берёт сначала правую (твоя задача)
//   B) Арбитр: официант разрешает есть не более N-1 одновременно
//   C) Chandy-Misra: вилки передаются через сообщения
//
// Реализуй вариант A или B. Объясни почему он не дедлочится.
//
// Запуск:
//   go run main.go
//   go run -race main.go   ← не должно быть ошибок

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const numPhilosophers = 5

type Fork struct {
	sync.Mutex
}

type Philosopher struct {
	id         int
	leftFork   *Fork
	rightFork  *Fork
	timesEaten int
}

// TODO: реализуй метод eat
// Вариант A: философ с нечётным id берёт сначала правую, потом левую вилку
// Вариант B: используй "официанта" (семафор на 4 разрешения)
func (p *Philosopher) eat(wg *sync.WaitGroup, stop <-chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			return
		default:
		}

		// Думаем
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

		// TODO: берём вилки (без дедлока!)
		// Подсказка: дедлок возникает когда все хватают вилки в одном порядке

		// Едим
		p.timesEaten++
		fmt.Printf("Философ %d ест (раз %d)\n", p.id, p.timesEaten)
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

		// TODO: кладём вилки
	}
}

func main() {
	forks := make([]*Fork, numPhilosophers)
	for i := range forks {
		forks[i] = &Fork{}
	}

	philosophers := make([]*Philosopher, numPhilosophers)
	for i := range philosophers {
		philosophers[i] = &Philosopher{
			id:        i,
			leftFork:  forks[i],
			rightFork: forks[(i+1)%numPhilosophers],
		}
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup

	for _, p := range philosophers {
		wg.Add(1)
		go p.eat(&wg, stop)
	}

	time.Sleep(500 * time.Millisecond)
	close(stop)
	wg.Wait()

	for _, p := range philosophers {
		fmt.Printf("Философ %d поел %d раз\n", p.id, p.timesEaten)
	}
}
