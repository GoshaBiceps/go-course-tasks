// Задача: Деление с проверкой
//
// Ожидаемый вывод:
//   result: 5
//   error: division by zero

package main

import (
"fmt"
"errors"
)


// TODO: напиши функцию safeDivide(a, b int) (int, error)
//       если b == 0, верни 0 и ошибку "division by zero"
//       иначе верни a / b и nil
func safeDivide(a, b int) (int, error) {
	if b == 0{
		return 0, errors.New("division by zero")
	}

	return a / b , nil 
}

func main() {
	// TODO: вызови safeDivide(10, 2) и выведи "result: <значение>"
	res, err := safeDivide(10, 2)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("result:", res)
	}
	// TODO: вызови safeDivide(10, 0) и выведи "error: <ошибка>"
	res, err = safeDivide(10, 0)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("result:", res)
	}
}
