// ============================================================
// Задача: Найди и исправь дедлоки  🔴 Senior
// ============================================================
//
// Вопрос-ловушка на собесах уровня Senior.
// "Что не так в этом коде?"
//
// Здесь 4 независимых сценария с дедлоками. Найди каждый и исправь.
//
// Запуск:
//   go run main.go           ← укажи номер сценария: go run main.go 1
//   go run -race main.go 4   ← для сценария с гонкой

package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// ============================================================
// Сценарий 1
// ============================================================
// TODO: запусти, пойми что происходит, исправь

func scenario1() {
	var mu1, mu2 sync.Mutex
	done := make(chan struct{})

	go func() {
		mu1.Lock()
		time.Sleep(1 * time.Millisecond)
		mu2.Lock()
		fmt.Println("A: захватил оба мьютекса")
		mu2.Unlock()
		mu1.Unlock()
		close(done)
	}()

	mu1.Lock()
	time.Sleep(1 * time.Millisecond)
	mu2.Lock()
	fmt.Println("B: захватил оба мьютекса")
	mu2.Unlock()
	mu1.Unlock()

	<-done
}

// ============================================================
// Сценарий 2
// ============================================================
// TODO: запусти, пойми что происходит, исправь

func scenario2() {
	ch := make(chan int)

	go func() {
		fmt.Println("sent:", <-ch)
	}()

	ch <- 42
}

// ============================================================
// Сценарий 3
// ============================================================
// TODO: запусти, пойми что происходит, исправь

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) incLocked() {
	c.count++
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.incLocked()
}

func (c *SafeCounter) IncAndLog() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.incLocked()
	c.incLocked()
}

func scenario3() {
	c := &SafeCounter{}
	c.IncAndLog()
	fmt.Println("count:", c.count)
}

// ============================================================
// Сценарий 4
// ============================================================
// TODO: запусти (с -race), пойми что происходит, исправь

func scenario4() {
	var wg sync.WaitGroup
	results := make([]int, 5)

	for i := 0; i < 5; i++ {
		n := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			results[n] = n * n
		}()
	}

	wg.Wait()
	fmt.Println(results)
}

func main() {
	scenario := 1
	if len(os.Args) > 1 {
		scenario, _ = strconv.Atoi(os.Args[1])
	}

	fmt.Printf("=== Сценарий %d ===\n", scenario)
	fmt.Println("В этом коде есть баг. Запусти, изучи поведение (трейсбек / -race), найди и исправь.")
	fmt.Println()

	switch scenario {
	case 1:
		scenario1()
	case 2:
		scenario2()
	case 3:
		scenario3()
	case 4:
		scenario4()
	default:
		fmt.Println("Укажи номер сценария: go run main.go [1-4]")
	}
}
