// Задача 3: Сумма чисел
//
// Ожидаемый вывод:
//   sum int: 15
//   sum float: 7.5

package main

import "fmt"

// TODO: объяви ограничение Number для ~int | ~float64
type Number interface {
	~int | ~float64
}
// TODO: напиши функцию Sum[T Number](items []T) T
//       возвращает сумму всех элементов среза
func Sum[T Number](items []T) T{
	var sum T 

	for _, num := range items {
		sum += num 
	}
	return sum 
}

func main() {
	 
	// TODO: вызови Sum для []int{1, 2, 3, 4, 5} и выведи "sum int: <результат>"
	 intSum := Sum([]int{1, 2, 3, 4, 5})
	 fmt.Printf("sum:%d\n", intSum)
	// TODO: вызови Sum для []float64{1.5, 2.0, 4.0} и выведи "sum float: <результат>"
	floatSum := Sum([]float64{1.5, 2.0, 4.0})
	fmt.Printf("sum:%f\n", floatSum)
}
