package main

import "fmt"

// TODO: напиши функцию First[T any](items []T) (T, bool)
//       если срез пустой — верни нулевое значение T и false
//       иначе — верни items[0] и true
func First[T any](items []T) (T, bool) {
	var zero T
	if len(items) == 0 {
		return zero, false
	}
	return items[0], true
}

func main() {
	// TODO: вызови First для []int{10, 20, 30} и выведи результат
	//       fmt.Printf("first: %v, ok: %v\n", val, ok)
	val, ok := First([]int{10, 20, 30})
	fmt.Printf("first: %v, ok: %v\n", val, ok)

	// TODO: вызови First для пустого []int{} и выведи результат
	val, ok = First([]int{})
	fmt.Printf("first: %v, ok: %v\n", val, ok)
}