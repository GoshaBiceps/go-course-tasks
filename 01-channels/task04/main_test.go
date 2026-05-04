package main

import (
	"context"
	"testing"
	"time"
)

func TestFastestReturnsFirst(t *testing.T) {
	ctx := context.Background()
	urls := []string{"a", "b", "c"}

	result, err := fastest(ctx, urls)
	if err != nil {
		t.Fatal("не ожидали ошибки:", err)
	}
	if result.URL == "" {
		t.Error("URL не должен быть пустым")
	}
}

func TestFastestTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Все запросы занимают >= 50ms, контекст истечёт раньше
	urls := []string{"slow1", "slow2"}
	_, err := fastest(ctx, urls)
	if err == nil {
		t.Fatal("ожидали ошибку таймаута")
	}
}

func TestWithTimeout_OK(t *testing.T) {
	result, err := withTimeout(500*time.Millisecond, func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "done", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != "done" {
		t.Errorf("got %q, want done", result)
	}
}

func TestWithTimeout_Expired(t *testing.T) {
	_, err := withTimeout(10*time.Millisecond, func() (string, error) {
		time.Sleep(500 * time.Millisecond)
		return "done", nil
	})
	if err == nil {
		t.Fatal("ожидали ошибку таймаута")
	}
}
