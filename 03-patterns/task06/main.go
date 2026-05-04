// ============================================================
// Задача: Errgroup — группа горутин с общей ошибкой  🟡 Middle
// ============================================================
//
// Реализуй аналог golang.org/x/sync/errgroup (не используя его).
//
//   type Group struct { ... }
//
//   func WithContext(ctx context.Context) (*Group, context.Context)
//   func (g *Group) Go(fn func() error)
//   func (g *Group) Wait() error
//   func (g *Group) SetLimit(n int)   // ограничение на число параллельных Go
//
// Семантика:
//   - Первая вернувшаяся ошибка — отменяет производный контекст
//   - Wait блокируется до завершения всех Go-функций
//   - Wait возвращает ПЕРВУЮ не-nil ошибку (или nil)
//   - SetLimit(n): Go блокируется пока запущено >= n горутин
//   - Если в Go-функции случилась паника — Wait должен её пробросить (бонус)
//
// Зачем: универсальный примитив для параллельного fan-out с общей отменой.
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Group struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
	cancel  context.CancelFunc
	sem     chan struct{} // nil если лимита нет
}

// TODO: реализуй WithContext
// Подсказка: нужен производный context.WithCancel; cancel вызывается при ПЕРВОЙ ошибке
func WithContext(ctx context.Context) (*Group, context.Context) {
	// TODO
	return nil, ctx
}

// TODO: реализуй Go
// Подсказка: учитывай лимит (sem) — если задан, он ограничивает число параллельных вызовов
// При ошибке — запомни первую (errOnce) и отмени ctx
func (g *Group) Go(fn func() error) {
	// TODO
}

// TODO: реализуй Wait
func (g *Group) Wait() error {
	// TODO
	return g.err
}

// TODO: реализуй SetLimit — после вызова Go ограничивает N параллельных
// Если задать 0 или отрицательное — сброс лимита
func (g *Group) SetLimit(n int) {
	// TODO
}

func main() {
	ctx := context.Background()
	g, ctx := WithContext(ctx)

	urls := []string{"a", "b", "c", "d"}

	for _, u := range urls {
		url := u
		g.Go(func() error {
			select {
			case <-time.After(50 * time.Millisecond):
				if url == "b" {
					return fmt.Errorf("fail on %s", url)
				}
				fmt.Printf("ok %s\n", url)
				return nil
			case <-ctx.Done():
				fmt.Printf("cancelled %s\n", url)
				return ctx.Err()
			}
		})
	}

	err := g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Println("первая ошибка:", err)
	}
}
