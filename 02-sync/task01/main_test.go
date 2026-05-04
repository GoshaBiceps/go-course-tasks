package main

import (
	"sync"
	"testing"
	"time"
)

func TestTTLCacheBasic(t *testing.T) {
	c := NewTTLCache[string, int](100 * time.Millisecond)

	c.Set("x", 42)
	v, ok := c.Get("x")
	if !ok || v != 42 {
		t.Errorf("Get(x) = %d, %v; want 42, true", v, ok)
	}
}

func TestTTLCacheExpiry(t *testing.T) {
	c := NewTTLCache[string, int](50 * time.Millisecond)
	c.Set("x", 42)
	time.Sleep(100 * time.Millisecond)

	_, ok := c.Get("x")
	if ok {
		t.Error("запись должна была устареть")
	}
}

func TestTTLCacheLen(t *testing.T) {
	c := NewTTLCache[string, int](100 * time.Millisecond)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	if c.Len() != 3 {
		t.Errorf("Len = %d, want 3", c.Len())
	}

	time.Sleep(150 * time.Millisecond)
	if c.Len() != 0 {
		t.Errorf("Len после TTL = %d, want 0", c.Len())
	}
}

func TestTTLCacheConcurrent(t *testing.T) {
	c := NewTTLCache[int, int](1 * time.Second)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Set(n, n*n)
			c.Get(n)
		}(i)
	}
	wg.Wait()
}
