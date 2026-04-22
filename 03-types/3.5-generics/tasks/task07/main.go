// Задача 7: Интеграционная задача — обобщенный Store
//
// Ожидаемый вывод:
//   strings: [go rust python]
//   ints: [1 2 3]
//   contains "go": true
//   contains "java": false
//   contains 2: true
//   contains 5: false

package main

import "fmt"

// TODO: объяви структуру Store[T any] с полем items []T
type Store[T any] struct {
	items []T
}

// TODO: добавь метод Add(item T) — добавляет элемент в items
func (s *Store[T]) Add(item T) {
	s.items = append(s.items, item)
}

// TODO: добавь метод All() []T — возвращает items
func (s *Store[T]) All() []T {
	result := make([]T, len(s.items))
	copy(result, s.items)
	return result //  возвращаем копию
}

// TODO: напиши функцию Contains[T comparable](items []T, target T) bool
//
//	возвращает true если target есть в items
func Contains[T comparable](items []T, target T) bool {
	for _, v := range items {
		if v == target {
			return true
		}
	}

	return false
}
func main() {
	// TODO: создай Store[string], добавь "go", "rust", "python"
	//       выведи "strings:", ss.All()
	ss := Store[string]{}
	ss.Add("go")
	ss.Add("rust")
	ss.Add("python")
	fmt.Println("strings:", ss.All())
	// TODO: создай Store[int], добавь 1, 2, 3
	//       выведи "ints:", si.All()
	si := Store[int]{}
	si.Add(1)
	si.Add(2)
	si.Add(3)
	fmt.Println("ints:", si.All())
	// TODO: выведи результаты Contains:
	//       Contains(ss.All(), "go")  → "contains \"go\": true"
	//       Contains(ss.All(), "java") → "contains \"java\": false"
	//       Contains(si.All(), 2)     → "contains 2: true"
	//       Contains(si.All(), 5)     → "contains 5: false"
	fmt.Println(`contains "go":`, Contains(ss.All(), "go"))
	fmt.Println(`contains "java":`, Contains(ss.All(), "java"))
	fmt.Println("contains 2:", Contains(si.All(), 2))
	fmt.Println("contains 5:", Contains(si.All(), 5))
}
