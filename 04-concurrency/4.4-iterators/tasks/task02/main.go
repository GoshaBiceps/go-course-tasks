// Задание 2: Бесконечный генератор чисел Фибоначчи
//
// Числа Фибоначчи: каждое следующее = сумма двух предыдущих.
// 0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, ...
//
// Напиши функцию Fibonacci() iter.Seq[int], которая генерирует
// числа Фибоначчи бесконечно.
//
// Внутри итератора:
//   a, b := 0, 1
//   для каждой итерации: yield(a), потом a, b = b, a+b
//
// В main() выведи первые 10 чисел и прерви цикл.
//
// Ожидаемый вывод:
//   0 1 1 2 3 5 8 13 21 34
//
// Запусти: go run main.go

package main

import (
	"fmt"
	"iter"
)

// TODO: напиши функцию Fibonacci() iter.Seq[int]
// Подсказка: используй две переменные a и b
// Не забудь проверить возвращаемое значение yield!
func Fibonacci() iter.Seq[int] {

	return func(yield func(int) bool) {
		a, b := 0, 1

		for {
			if !yield(a) {
				return
			}
			a, b = b, a+b
		}
	}
}

func main() {
	count := 0
	for n := range Fibonacci() {
		fmt.Print(n, " ")
		count++
		if count == 10 {
			break // прерываем бесконечный итератор
		}
	}
	fmt.Println()

	_ = iter.Seq[int](nil) // убери когда начнёшь использовать
}
