// Задача 6: Значения map в срез
//
// Ожидаемый вывод (порядок может отличаться):
//   values: [1 2 3]

package main

import "fmt"

// TODO: напиши функцию Values[K comparable, V any](m map[K]V) []V
//       возвращает срез всех значений из map
func Values[K comparable, V any] (m map[K]V) []V {
	var result []V 

	for _, v := range m {
		result = append(result, v)
	}
	
	return result 
}
func main() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	// TODO: вызови Values(m) и выведи результат через fmt.Println("values:", vals)
	 vals := Values(m)
	 fmt.Println("values:", vals)

}
