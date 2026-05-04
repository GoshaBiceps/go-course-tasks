package main

import (
	"runtime"
	"sort"
	"testing"
)

func TestFanIn(t *testing.T) {
	ch1 := make(chan int, 3)
	ch2 := make(chan int, 3)
	ch1 <- 1
	ch1 <- 3
	close(ch1)
	ch2 <- 2
	ch2 <- 4
	close(ch2)

	var got []int
	for v := range fanIn(ch1, ch2) {
		got = append(got, v)
	}
	sort.Ints(got)

	if len(got) != 4 {
		t.Fatalf("ожидали 4 элемента, получили %d: %v", len(got), got)
	}
	want := []int{1, 2, 3, 4}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d]: got %d, want %d", i, got[i], want[i])
		}
	}
}

func TestFanOutFanIn_NoLeaks(t *testing.T) {
	goroutinesBefore := runtime.NumGoroutine()

	source := make(chan int, 20)
	for i := 1; i <= 20; i++ {
		source <- i
	}
	close(source)

	workers := fanOut(source, 4)
	var processed []<-chan int
	for _, w := range workers {
		processed = append(processed, process(w))
	}

	sum := 0
	for v := range fanIn(processed...) {
		sum += v
	}

	// (1+2+...+20)*2 = 420
	if sum != 420 {
		t.Errorf("sum = %d, want 420", sum)
	}

	// Проверяем отсутствие утечек горутин
	goroutinesAfter := runtime.NumGoroutine()
	if goroutinesAfter > goroutinesBefore+2 {
		t.Errorf("возможная утечка горутин: до=%d, после=%d", goroutinesBefore, goroutinesAfter)
	}
}
