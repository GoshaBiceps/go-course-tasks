// Задача 5: Безопасный type assertion
//
// Ожидаемый вывод:
//   length: 5
//   not a string

package main

import "fmt"

// TODO: напиши функцию printStringLength(x any)
//       используй безопасный type assertion (двузначная форма)
//       если x — string, выводи "length: <len>"
//       иначе выводи "not a string"

func printStringLength(x any) {
	s, ok := x.(string)
	if ok {
		fmt.Println("lenght:", len(s))
		return 
	}
	
	fmt.Println("not a string")
}

func main() {
	printStringLength("hello")
	printStringLength(42)
}
