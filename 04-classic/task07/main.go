// ============================================================
// Задача: Readers-Writers без голодания  🔴 Senior
// ============================================================
//
// Классика на собесах. Отличается от задачи "WriterPriorityRWMutex" тем
// что здесь надо реализовать ДВА разных варианта и сравнить поведение.
//
// 1. Readers-preferring:
//    Читатели проходят всегда когда возможно. Писатель может голодать.
//
// 2. Fair (FIFO):
//    Порядок прихода соблюдается: если писатель пришёл раньше последующих
//    читателей — он проходит первым.
//
// Реализуй оба варианта с одинаковым интерфейсом:
//
//   type Lock interface {
//       RLock()
//       RUnlock()
//       Lock()
//       Unlock()
//   }
//
//   func NewReaderPreferring() Lock
//   func NewFair() Lock
//
// Задача не в коде — а в понимании семантики.
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

type Lock interface {
	RLock()
	RUnlock()
	Lock()
	Unlock()
}

// === Reader-preferring ===

type readerPref struct {
	// TODO: какие поля нужны?
	mu sync.Mutex
}

// TODO: реализуй NewReaderPreferring
func NewReaderPreferring() Lock {
	return &readerPref{}
}

// TODO
func (l *readerPref) RLock()   {}
func (l *readerPref) RUnlock() {}

// TODO
func (l *readerPref) Lock()   {}
func (l *readerPref) Unlock() {}

// === Fair (FIFO) ===

type fair struct {
	// TODO: подумай про очередь ожидающих — каждый "в очереди" со своим сигналом
	mu sync.Mutex
}

// TODO: реализуй NewFair
func NewFair() Lock {
	return &fair{}
}

// TODO
func (l *fair) RLock()   {}
func (l *fair) RUnlock() {}

// TODO
func (l *fair) Lock()   {}
func (l *fair) Unlock() {}

// === Демо ===

func demo(name string, l Lock) {
	fmt.Printf("\n=== %s ===\n", name)
	var reads, writes atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 5 {
				l.RLock()
				reads.Add(1)
				time.Sleep(3 * time.Millisecond)
				l.RUnlock()
			}
		}()
	}

	// Писатель посреди потока читателей
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond)
		t0 := time.Now()
		l.Lock()
		fmt.Printf("писатель зашёл через %v\n", time.Since(t0))
		time.Sleep(10 * time.Millisecond)
		writes.Add(1)
		l.Unlock()
	}()

	wg.Wait()
	fmt.Printf("reads=%d writes=%d\n", reads.Load(), writes.Load())
}

func main() {
	demo("reader-preferring", NewReaderPreferring())
	demo("fair", NewFair())
}
