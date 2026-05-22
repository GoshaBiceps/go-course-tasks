// ============================================================
// Задача: Token Bucket Rate Limiter  🟡 Middle
// ============================================================
//
// Вопрос с собесов уровня Middle.
//
// Реализуй rate limiter без сторонних библиотек (не golang.org/x/time/rate).
//
// Алгоритм Token Bucket:
//   - Есть "ведро" ёмкостью capacity токенов
//   - Каждые 1/rate секунд добавляется 1 токен (но не больше capacity)
//   - Allow() забирает 1 токен и возвращает true, или false если ведро пусто
//
// Реализуй ДВА варианта:
//
//   1. TokenBucket — на основе time.Ticker и горутины
//   2. LazyTokenBucket — ленивый: считает токены математически при каждом вызове
//      (без горутин! Используй time.Since для вычисления накопленных токенов)
//
// LazyTokenBucket предпочтителен в продакшне: не создаёт горутины.
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// === Вариант 1: с горутиной ===

type TokenBucket struct {
	tokens   atomic.Int64
	capacity int64
	quit     chan struct{}
}

// TODO: реализуй NewTokenBucket
// Подсказка: фоновая горутина добавляет токены с нужной частотой; начинай с полным ведром
func NewTokenBucket(rate float64, capacity int64) *TokenBucket {

	tb := &TokenBucket{
		capacity: capacity,
		quit:     make(chan struct{}),
	}

	tb.tokens.Store(capacity) // стартуем сразу с токенами

	interval := time.Duration(float64(time.Second) / rate)
	ticker := time.NewTicker(interval) // создали структуру тикер , у нее есть поле с каналом  C

	go func() {
		for {
			select {
			case <-ticker.C: // канал таймера; каждые interval  отправляет сигнал .

				current := tb.tokens.Load()

				if current < tb.capacity {
					tb.tokens.Add(1)
				}
			case <-tb.quit:
				return
			}
		}

	}()
	return tb
}

// TODO: Allow забирает 1 токен. Возвращает false если ведро пусто.
// Подсказка: операция должна быть потокобезопасной без мьютекса
func (tb *TokenBucket) Allow() bool {

	for {
		currency := tb.tokens.Load()

		if currency <= 0 {
			return false
		}

		if tb.tokens.CompareAndSwap(currency, currency-1) {
			return true
		}

	}
}

func (tb *TokenBucket) Close() { close(tb.quit) }

// === Вариант 2: ленивый (без горутин) ===

type LazyTokenBucket struct {
	mu         sync.Mutex
	tokens     float64
	capacity   float64
	rate       float64 // токенов в секунду
	lastRefill time.Time
}

// TODO: реализуй NewLazyTokenBucket
func NewLazyTokenBucket(rate, capacity float64) *LazyTokenBucket {

	return &LazyTokenBucket{
		tokens:     capacity,   // текущее количество токенов
		capacity:   capacity,   //максимальный размер бакета
		rate:       rate,       // токенов в секунду
		lastRefill: time.Now(), // когда последний раз подсчитывали бакет
	}
}

// TODO: реализуй Allow для LazyTokenBucket
// Подсказка: горутина не нужна — используй lastRefill и time.Since между вызовами
func (lb *LazyTokenBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	elapsed := time.Since(lb.lastRefill).Seconds() //  считет сколько прошло секунд с полседнего  refil

	refill := elapsed * lb.rate                    //   вычисляем сколько токенов накопилось
	lb.tokens = min(lb.capacity, lb.tokens+refill) // добавляем токены
	lb.lastRefill = time.Now()

	if lb.tokens < 1 {
		return false
	}

	lb.tokens--

	return true
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("=== TokenBucket (10 req/s, burst 3) ===")
	tb := NewTokenBucket(10, 3)
	defer tb.Close()

	for i := 0; i < 6; i++ {
		ok := tb.Allow()
		fmt.Printf("запрос %d: %v\n", i+1, ok)
		if i == 2 {
			time.Sleep(200 * time.Millisecond) // ждём накопления токенов
		}
	}

	fmt.Println("\n=== LazyTokenBucket (5 req/s, burst 2) ===")
	lb := NewLazyTokenBucket(5, 2)

	for i := 0; i < 5; i++ {
		ok := lb.Allow()
		fmt.Printf("запрос %d: %v\n", i+1, ok)
		time.Sleep(100 * time.Millisecond)
	}
}
