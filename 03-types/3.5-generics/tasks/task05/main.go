// Задача 5: Обобщенная структура Pair
//
// Ожидаемый вывод:
//   name: Alice
//   count: 42

package main

import "fmt"

// TODO: объяви структуру Pair[K comparable, V any] с полями Key K, Value V
type Pair[K comparable, V any] struct{
	Key K
	Value V 
}

func main() {
	p1 := Pair[string, string]{Key: "name", Value: "Alice"}
	fmt.Printf("%s: %s\n", p1.Key, p1.Value)

	p2 := Pair[string, int]{Key: "count", Value: 42}
	fmt.Printf("%s: %d\n", p2.Key, p2.Value)
}
