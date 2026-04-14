// Задача 4: Функция с указателем на базовый тип
//
// Ожидаемый вывод:
//   score: 0

package main

import "fmt"

// TODO: напиши функцию resetScore(score *int)
//       она должна устанавливать значение по указателю в 0
func resetScore(score *int){
	*score = 0 
}
func main() {
	score := 42

	// TODO: вызови resetScore передав &score
	resetScore(&score)
	
	fmt.Println("score:", score)
}
