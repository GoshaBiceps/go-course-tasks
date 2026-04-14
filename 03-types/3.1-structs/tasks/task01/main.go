// Задание 1: Профиль пользователя
//
// Тебе нужно:
// 1. Создать структуру Profile с полями Name (string), Age (int), IsActive (bool)
// 2. В main создать значение через именованный литерал
// 3. Вывести поля в нужном формате
//
// Ожидаемый вывод:
//   Name: Anna, Age: 25, IsActive: true
//
// Запусти: go run main.go

package main

import "fmt"

type Profile struct {
	Name string
	Age int 
	IsActive bool 
}

// TODO: объяви структуру Profile с полями Name, Age, IsActive

func main() {
	// TODO: создай значение Profile с Name: "Anna", Age: 25, IsActive: true
	p := Profile{
		Name: "Anna",
		Age: 25,
		IsActive: true,
	}

	fmt.Printf("Name: %s, Age: %d, IsActive: %t\n", p.Name, p.Age, p.IsActive)
}