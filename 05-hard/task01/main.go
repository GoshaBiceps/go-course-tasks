// ============================================================
// Задача: Шардированная конкурентная map  🔴 Senior
// ============================================================
//
// Вопрос с собесов уровня Senior+.
//
// Проблема sync.Map: при высокой нагрузке один мьютекс становится узким местом.
// Решение: шардирование — делим map на N независимых секций с отдельными мьютексами.
//
//   type ShardedMap[K comparable, V any] struct { ... }
//
//   func NewShardedMap[K comparable, V any](shards int) *ShardedMap[K, V]
//   func (m *ShardedMap[K, V]) Set(key K, value V)
//   func (m *ShardedMap[K, V]) Get(key K) (V, bool)
//   func (m *ShardedMap[K, V]) Delete(key K)
//   func (m *ShardedMap[K, V]) Len() int
//   func (m *ShardedMap[K, V]) Range(fn func(K, V) bool)
//
// Номер шарда определяется хешом ключа:
//   shardIndex = hash(key) % numShards
//
// Бенчмарк: сравни с sync.Map при 90% чтений / 10% записей.
//
// Проверь:
//   go test -race -v ./...
//   go test -bench=. -benchmem ./...

package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type shard[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

type ShardedMap[K comparable, V any] struct {
	shards []*shard[K, V]
	n      int
}

// TODO: реализуй NewShardedMap
func NewShardedMap[K comparable, V any](numShards int) *ShardedMap[K, V] {
	m := &ShardedMap[K, V]{
		shards: make([]*shard[K, V], numShards),
		n:      numShards,
	}
	for i := range m.shards {
		m.shards[i] = &shard[K, V]{items: make(map[K]V)}
	}
	return m
}

// shardFor возвращает шард для данного ключа
func (m *ShardedMap[K, V]) shardFor(key K) *shard[K, V] {
	h := fnv.New32a()
	fmt.Fprintf(h, "%v", key)
	return m.shards[h.Sum32()%uint32(m.n)]
}

// TODO: реализуй Set
func (m *ShardedMap[K, V]) Set(key K, value V) {
}

// TODO: реализуй Get — подумай какой тип блокировки уместен для чтения
func (m *ShardedMap[K, V]) Get(key K) (V, bool) {
	var zero V
	return zero, false
}

// TODO: реализуй Delete
func (m *ShardedMap[K, V]) Delete(key K) {
}

// TODO: реализуй Len — сумма размеров всех шардов
func (m *ShardedMap[K, V]) Len() int {
	return 0
}

// TODO: реализуй Range — обходит все элементы всех шардов
// fn возвращает false — прерывает обход
// Подсказка: не держи мьютекс во время вызова fn (fn может быть медленной или вызывать Set/Delete на той же map)
func (m *ShardedMap[K, V]) Range(fn func(K, V) bool) {
	// TODO
}

func main() {
	m := NewShardedMap[string, int](16)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		n := i
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("key%d", n)
			m.Set(key, n)
			v, _ := m.Get(key)
			_ = v
		}()
	}
	wg.Wait()

	fmt.Printf("Len = %d\n", m.Len())
}
