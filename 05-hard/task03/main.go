// ============================================================
// Задача: Connection Pool  🔴 Senior
// ============================================================
//
// Вопрос с собесов уровня Senior.
//
// Реализуй пул соединений с базой данных:
//
//   type Pool struct { ... }
//
//   func NewPool(maxConn int, factory func() (Conn, error)) *Pool
//   func (p *Pool) Acquire(ctx context.Context) (Conn, error)
//   func (p *Pool) Release(conn Conn)
//   func (p *Pool) Close()
//   func (p *Pool) Stats() PoolStats
//
// Требования:
//   - Не более maxConn одновременных соединений
//   - Acquire блокируется пока нет свободного соединения
//   - Если ctx отменён во время ожидания — вернуть ctx.Err()
//   - Соединения переиспользуются (не создаём новое на каждый Acquire)
//   - Health check: если соединение "сломано" — создаём новое вместо него
//   - Close закрывает все незанятые соединения, ждёт возврата занятых
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var ErrPoolClosed = errors.New("пул закрыт")
var ErrPoolExhausted = errors.New("пул исчерпан")

type Conn interface {
	Ping() error
	Close() error
	ID() int
}

type PoolStats struct {
	Total    int
	Idle     int
	InUse    int
	Acquired int64
	Released int64
}

type mockConn struct {
	id     int
	broken bool
}

func (c *mockConn) Ping() error {
	if c.broken {
		return errors.New("соединение сломано")
	}
	return nil
}
func (c *mockConn) Close() error {
	fmt.Printf("закрыто соединение %d\n", c.id)
	return nil
}
func (c *mockConn) ID() int { return c.id }

var connIDCounter atomic.Int32

type Pool struct {
	mu       sync.Mutex
	cond     *sync.Cond
	idle     []Conn
	inUse    int
	maxConn  int
	factory  func() (Conn, error)
	closed   bool
	acquired atomic.Int64
	released atomic.Int64
}

// TODO: реализуй NewPool
func NewPool(maxConn int, factory func() (Conn, error)) *Pool {
	p := &Pool{
		maxConn: maxConn,
		factory: factory,
	}
	p.cond = sync.NewCond(&p.mu)
	return p
}

// TODO: реализуй Acquire
// Подсказка: три сценария: idle есть, можно создать, нужно ждать
// Отмена ctx должна разбудить ожидающего — подумай как
func (p *Pool) Acquire(ctx context.Context) (Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		if p.closed { // проверка  можно выдать соединение или нет
			return nil, ErrPoolClosed
		}

		// 1. Есть свободное соединение
		if len(p.idle) > 0 {
			n := len(p.idle) - 1
			conn := p.idle[n]
			p.idle = p.idle[:n]

			// health check
			if conn.Ping() != nil {
				conn.Close()
				continue
			}

			p.inUse++
			p.acquired.Add(1)
			return conn, nil
		}

		// 2. Можно создать новое соединение
		if p.inUse < p.maxConn {
			p.inUse++

			p.mu.Unlock()
			conn, err := p.factory()
			p.mu.Lock()

			if err != nil {
				p.inUse--
				p.cond.Signal()
				return nil, err
			}

			p.acquired.Add(1)
			return conn, nil
		}

		// 3. Всё занято — ждём
		done := make(chan struct{})

		go func() {
			select {
			case <-ctx.Done():
				p.mu.Lock()
				p.cond.Broadcast()
				p.mu.Unlock()
			case <-done:
			}
		}()

		p.cond.Wait()
		close(done)

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}
}

// TODO: реализуй Release
// Подсказка: после возврата нужно разбудить ожидающего
func (p *Pool) Release(conn Conn) {

	p.mu.Lock()
	defer p.mu.Unlock()

	// pool уже закрыт
	if p.closed {

		// conn больше не нужен
		conn.Close()

		p.inUse--

		p.released.Add(1)

		// будим waiting goroutine
		p.cond.Broadcast()

		return
	}

	// conn снова становится idle
	p.idle = append(p.idle, conn)

	// conn больше не используется
	p.inUse--

	// статистика release
	p.released.Add(1)

	// будим одну waiting goroutine
	p.cond.Signal()
}

// TODO: реализуй Close
func (p *Pool) Close() {

	p.mu.Lock()

	// уже закрыт
	if p.closed {
		p.mu.Unlock()
		return
	}

	// помечаем pool как closed
	p.closed = true

	// закрываем все idle conn
	for _, conn := range p.idle {
		conn.Close()
	}

	// idle больше нет
	p.idle = nil

	// будим все waiting goroutine
	p.cond.Broadcast()

	// ждём пока вернутся все inUse conn
	for p.inUse > 0 {

		// Wait:
		// - отпускает mutex
		// - goroutine засыпает
		// - после Signal/Broadcast
		//   снова захватывает mutex
		p.cond.Wait()
	}

	p.mu.Unlock()
}

func (p *Pool) Stats() PoolStats {
	p.mu.Lock()
	defer p.mu.Unlock()
	return PoolStats{
		Total:    p.inUse + len(p.idle),
		Idle:     len(p.idle),
		InUse:    p.inUse,
		Acquired: p.acquired.Load(),
		Released: p.released.Load(),
	}
}

func main() {
	factory := func() (Conn, error) {
		id := int(connIDCounter.Add(1))
		fmt.Printf("создано соединение %d\n", id)
		return &mockConn{id: id}, nil
	}

	pool := NewPool(3, factory)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			conn, err := pool.Acquire(ctx)
			if err != nil {
				fmt.Printf("горутина %d: ошибка %v\n", n, err)
				return
			}
			fmt.Printf("горутина %d: соединение %d\n", n, conn.ID())
			time.Sleep(50 * time.Millisecond)
			pool.Release(conn)
		}(i)
	}

	wg.Wait()
	stats := pool.Stats()
	fmt.Printf("\nСтатистика: acquired=%d, released=%d, idle=%d\n",
		stats.Acquired, stats.Released, stats.Idle)
	pool.Close()
}
