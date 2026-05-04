// ============================================================
// Задача: Bounded Generator с Backpressure  🔴 Senior
// ============================================================
//
// Частый вопрос на собесах уровня Senior.
//
// Реализуй Generator — источник данных с ограниченным буфером.
// Если потребитель не успевает читать — Generator замедляется (backpressure).
//
// Интерфейс:
//
//   type Generator struct { ... }
//
//   func NewGenerator(bufSize int) *Generator
//   func (g *Generator) Send(v int) bool   // false если full и timeout
//   func (g *Generator) Chan() <-chan int
//   func (g *Generator) Close()
//
// Требования:
//   - Send блокируется максимум sendTimeout (100мс)
//   - Если за sendTimeout буфер не освободился — Send возвращает false
//   - Close завершает работу: Chan закрывается, Send возвращает false
//   - Безопасен для одновременного использования из нескольких горутин
//
// Дополнительно реализуй:
//   - Stats() (sent, dropped int64) — количество успешно отправленных и дропнутых
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

const sendTimeout = 100 * time.Millisecond

type Generator struct {
	ch      chan int
	closed  chan struct{}
	once    sync.Once
	sent    atomic.Int64
	dropped atomic.Int64
}

// TODO: реализуй NewGenerator
func NewGenerator(bufSize int) *Generator {
	// TODO: создай Generator с буферизованным ch и каналом-сигналом closed
	return nil
}

// TODO: реализуй Send
// Отправляет значение в буферизованный канал.
// Если буфер полон — ждёт до sendTimeout, потом возвращает false.
// Если Generator закрыт — сразу возвращает false.
func (g *Generator) Send(v int) bool {
	// TODO: три сценария: успешная отправка, таймаут, генератор закрыт — каждый должен обновлять статистику
	return false
}

// Chan возвращает канал для чтения данных
func (g *Generator) Chan() <-chan int {
	return g.ch
}

// TODO: реализуй Close — закрой closed канал через sync.Once, закрой ch
func (g *Generator) Close() {
	// TODO: закрытие должно быть безопасным при параллельных вызовах
}

// Stats возвращает статистику
func (g *Generator) Stats() (sent, dropped int64) {
	return g.sent.Load(), g.dropped.Load()
}

func main() {
	gen := NewGenerator(5)

	// Медленный потребитель
	go func() {
		for v := range gen.Chan() {
			fmt.Printf("получено: %d\n", v)
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Println("канал закрыт")
	}()

	// Быстрый производитель
	for i := 0; i < 20; i++ {
		ok := gen.Send(i)
		if !ok {
			fmt.Printf("дропнуто: %d\n", i)
		}
	}

	gen.Close()
	time.Sleep(500 * time.Millisecond)

	sent, dropped := gen.Stats()
	fmt.Printf("Отправлено: %d, Дропнуто: %d\n", sent, dropped)
}
