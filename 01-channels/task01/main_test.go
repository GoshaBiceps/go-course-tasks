package main

import (
	"testing"
)

func TestPipeline(t *testing.T) {
	got := collect(filterEven(square(generate(1, 2, 3, 4, 5))))
	want := []int{4, 16}

	if len(got) != len(want) {
		t.Fatalf("длина: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d]: got %d, want %d", i, got[i], want[i])
		}
	}
}

func TestPipelineEmpty(t *testing.T) {
	got := collect(filterEven(square(generate())))
	if len(got) != 0 {
		t.Errorf("ожидали пустой срез, получили %v", got)
	}
}

func TestPipelineAllOdd(t *testing.T) {
	got := collect(filterEven(square(generate(1, 3, 5))))
	// 1, 9, 25 — все нечётные
	if len(got) != 0 {
		t.Errorf("ожидали пустой срез (все нечётные квадраты), получили %v", got)
	}
}

func collect(ch <-chan int) []int {
	var result []int
	for v := range ch {
		result = append(result, v)
	}
	return result
}

func BenchmarkPipeline(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i
	}
	for b.Loop() {
		for range filterEven(square(generate(nums...))) {
		}
	}
}
