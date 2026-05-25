// ============================================================
// Задача: Потокобезопасный LRU-кеш  🔴 Senior
// ============================================================
//
// Один из самых популярных вопросов на собесах — LeetCode 146 + конкурентность.
//
// Реализуй LRU кеш с:
//   - O(1) Get и Put
//   - Ограничением capacity: при переполнении удаляется давно неиспользованный элемент
//   - Потокобезопасностью (sync.RWMutex)
//
//   type LRUCache[K comparable, V any] struct { ... }
//
//   func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V]
//   func (c *LRUCache[K, V]) Get(key K) (V, bool)
//   func (c *LRUCache[K, V]) Put(key K, value V)
//   func (c *LRUCache[K, V]) Len() int
//
// Алгоритм:
//   - doubly linked list: голова = самый недавний, хвост = самый старый
//   - map[K]*node для O(1) доступа
//   - Get перемещает узел в начало списка
//   - Put: если ключ есть — обновить и переместить в начало
//           если нет — добавить в начало, если len > cap — удалить хвостz
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"container/list"
	"fmt"
	"sync"
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	mu    sync.RWMutex
	cap   int
	list  *list.List
	items map[K]*list.Element
}

// TODO: реализуй NewLRUCache
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		cap:   capacity,
		list:  list.New(),
		items: make(map[K]*list.Element),
	}
}

// TODO: реализуй Get
// Подсказка: чтение тоже требует записи — подумай почему
func (c *LRUCache[K, V]) Get(key K) (V, bool) {

	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V

	elem, ok := c.items[key]
	if !ok {
		return zero, false
	}

	// двигаем элемент в начало списка
	c.list.MoveToFront(elem)

	ent := elem.Value.(entry[K, V])

	return ent.value, true
}

// TODO: реализуй Put
// Подсказка: если ключ уже есть — обнови и повысь в приоритете;
// если переполнено — вытесни least-recently-used
func (c *LRUCache[K, V]) Put(key K, value V) {

	c.mu.Lock()
	defer c.mu.Unlock()

	// ключ уже есть
	if elem, ok := c.items[key]; ok {

		// обновляем value
		elem.Value = entry[K, V]{
			key:   key,
			value: value,
		}

		// делаем элемент самым свежим
		c.list.MoveToFront(elem)

		return
	}

	// создаём новый элемент
	elem := c.list.PushFront(entry[K, V]{
		key:   key,
		value: value,
	})

	// сохраняем указатель на node списка
	c.items[key] = elem

	// проверяем переполнение cache
	if c.list.Len() > c.cap {

		// берём самый старый элемент
		last := c.list.Back()

		if last != nil {

			// получаем entry
			ent := last.Value.(entry[K, V])

			// удаляем из map
			delete(c.items, ent.key)

			// удаляем из linked list
			c.list.Remove(last)
		}
	}
}

func (c *LRUCache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.list.Len()
}

func main() {
	cache := NewLRUCache[int, string](3)

	cache.Put(1, "a")
	cache.Put(2, "b")
	cache.Put(3, "c")

	v, ok := cache.Get(1) // 1 теперь самый свежий
	fmt.Printf("Get(1) = %q, ok=%v\n", v, ok)

	cache.Put(4, "d") // вытесняет 2 (самый старый после доступа к 1)

	_, ok2 := cache.Get(2)
	fmt.Printf("Get(2) после вытеснения: ok=%v\n", ok2) // false

	v3, _ := cache.Get(3)
	v4, _ := cache.Get(4)
	fmt.Printf("Get(3) = %q, Get(4) = %q\n", v3, v4)
	fmt.Printf("Len = %d\n", cache.Len()) // 3
}
