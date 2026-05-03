// Задание 1: Итератор для стека
//
// Стек - это структура данных где последний добавленный элемент
// достаётся первым (как стопка тарелок).
//
// Реализуй:
//   type Stack struct { items []int }
//   func (s *Stack) Push(val int)      - добавить элемент наверх
//   func (s *Stack) All() iter.Seq[int] - итератор от вершины к основанию
//
// В main():
//   1. Создай стек и добавь числа 1, 2, 3, 4, 5
//   2. Выведи все элементы через for range
//   3. Выведи элементы пока не встретишь 3 (используй break)
//
// Ожидаемый вывод:
//   Все элементы стека:
//   5
//   4
//   3
//   2
//   1
//
//   До элемента 3:
//   5
//   4
//
// Запусти: go run main.go

package main

import (
	"fmt"
	"iter"
)

// TODO: напиши структуру Stack
type Stack struct {
	items []int
}

// TODO: напиши метод Push(val int)
func (s *Stack) Push(val int) {
	s.items = append(s.items, val)
}

// TODO: напиши метод All() iter.Seq[int]
// Подсказка: перебирай items с конца (от len-1 до 0)
// Не забудь проверять возвращаемое значение yield и прерваться если false
func (s *Stack) All() iter.Seq[int] {

	return func(yield func(int) bool) {
		for i := len(s.items) - 1; i >= 0; i-- {
			if !yield(s.items[i]) {
				return
			}
		}
	}

}
func main() {
	// TODO: создай стек, добавь элементы 1..5
	s := &Stack{}
	fmt.Println("Все элементы стека:")
	for i := 1; i <= 5; i++ {
		s.Push(i)
	}
	// TODO: for range stack.All() { ... }
	for v := range s.All() {
		fmt.Println(v)
	}
	fmt.Println("\nДо элемента 3:")
	// TODO: for range stack.All() { if v == 3 { break } ... }
	for v := range s.All() {
		if v == 3 {
			break
		}
		fmt.Println(v)
	}
}
