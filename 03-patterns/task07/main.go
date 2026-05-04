// ============================================================
// Задача: Circuit Breaker  🔴 Senior
// ============================================================
//
// Паттерн для защиты от лавинообразных падений при проблемах
// у внешнего сервиса.
//
// Три состояния:
//   Closed    — всё ок, запросы пропускаются, считаем последовательные ошибки
//   Open      — запросы моментально фейлятся (не идём в бэкенд) до таймаута
//   HalfOpen  — пробный запрос: если успешен — возвращаемся в Closed;
//               если падает — снова Open
//
// Интерфейс:
//
//   type CircuitBreaker struct { ... }
//
//   func New(failureThreshold int, openTimeout time.Duration) *CircuitBreaker
//   func (cb *CircuitBreaker) Call(fn func() error) error
//   func (cb *CircuitBreaker) State() State
//
// Требования:
//   - Потокобезопасен (много параллельных Call)
//   - В состоянии Open возвращает ErrOpen сразу, не вызывая fn
//   - В HalfOpen пропускает ТОЛЬКО один пробный Call одновременно
//   - Успех в HalfOpen → Closed, счётчик ошибок сбрасывается
//   - Ошибка в HalfOpen → снова Open, таймер открытия продлевается
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrOpen = errors.New("circuit breaker: открыт")

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	return [...]string{"Closed", "Open", "HalfOpen"}[s]
}

type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failures         int
	failureThreshold int
	openTimeout      time.Duration
	openedAt         time.Time
	probeInFlight    bool
}

// TODO: реализуй конструктор
func New(failureThreshold int, openTimeout time.Duration) *CircuitBreaker {
	return nil
}

// TODO: реализуй Call
// Подсказка: проверь state перед вызовом fn; после вызова — обнови состояние по результату
// Отдельная сложность — переход Open → HalfOpen по времени (не через таймер, а лениво при Call)
func (cb *CircuitBreaker) Call(fn func() error) error {
	// TODO
	return nil
}

// TODO: реализуй State — просто читаем текущее состояние (не забудь про ленивый переход Open → HalfOpen)
func (cb *CircuitBreaker) State() State {
	// TODO
	return StateClosed
}

func main() {
	cb := New(3, 100*time.Millisecond)

	failingCall := func() error { return errors.New("backend error") }
	goodCall := func() error { return nil }

	// Роняем до Open
	for i := 0; i < 5; i++ {
		err := cb.Call(failingCall)
		fmt.Printf("попытка %d: state=%s err=%v\n", i+1, cb.State(), err)
	}

	// В Open — мгновенный ErrOpen
	err := cb.Call(goodCall)
	fmt.Printf("в Open: state=%s err=%v\n", cb.State(), err)

	// Ждём openTimeout — следующий Call должен быть probe (HalfOpen)
	time.Sleep(150 * time.Millisecond)
	err = cb.Call(goodCall)
	fmt.Printf("probe: state=%s err=%v\n", cb.State(), err) // должен быть Closed, err=nil
}
