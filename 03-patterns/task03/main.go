// ============================================================
// Задача: Pub/Sub Broker  🔴 Senior
// ============================================================
//
// Вопрос с собесов уровня Senior.
//
// Реализуй in-memory брокер сообщений:
//
//   type Broker[T any] struct { ... }
//
//   func NewBroker[T any]() *Broker[T]
//   func (b *Broker[T]) Subscribe(topic string) <-chan T
//   func (b *Broker[T]) Unsubscribe(topic string, ch <-chan T)
//   func (b *Broker[T]) Publish(topic string, msg T)
//   func (b *Broker[T]) Close()
//
// Требования:
//   - Один топик может иметь несколько подписчиков
//   - Publish не блокируется — медленный подписчик дропает сообщения (буфер 10)
//   - Unsubscribe корректно удаляет подписчика и закрывает его канал
//   - Close завершает брокер: закрывает все каналы всех подписчиков
//   - Безопасен для параллельного использования
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"fmt"
	"sync"
	"time"
)

type Broker[T any] struct {
	mu          sync.RWMutex
	subscribers map[string][]chan T
	closed      bool
}

// TODO: реализуй NewBroker
func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		subscribers: make(map[string][]chan T),
	}
}

// TODO: реализуй Subscribe
// Подсказка: каждый подписчик — это отдельный канал; учти состояние closed
func (b *Broker[T]) Subscribe(topic string) <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()
	// TODO
	return nil
}

// TODO: реализуй Unsubscribe
func (b *Broker[T]) Unsubscribe(topic string, sub <-chan T) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// TODO
}

// TODO: реализуй Publish
// Подсказка: медленный подписчик не должен блокировать остальных
func (b *Broker[T]) Publish(topic string, msg T) {
	// TODO
}

// TODO: реализуй Close
func (b *Broker[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	// TODO
}

func main() {
	broker := NewBroker[string]()

	sub1 := broker.Subscribe("news")
	sub2 := broker.Subscribe("news")
	sub3 := broker.Subscribe("sports")

	// Подписчики читают в фоне
	printAll := func(name string, ch <-chan string) {
		for msg := range ch {
			fmt.Printf("[%s] %s\n", name, msg)
		}
		fmt.Printf("[%s] канал закрыт\n", name)
	}

	go printAll("sub1-news", sub1)
	go printAll("sub2-news", sub2)
	go printAll("sub3-sports", sub3)

	broker.Publish("news", "Статья 1")
	broker.Publish("news", "Статья 2")
	broker.Publish("sports", "Матч 1")

	// Отписываем sub2
	broker.Unsubscribe("news", sub2)

	broker.Publish("news", "Статья 3") // только sub1 получит

	time.Sleep(100 * time.Millisecond)
	broker.Close()
	time.Sleep(50 * time.Millisecond)
}
