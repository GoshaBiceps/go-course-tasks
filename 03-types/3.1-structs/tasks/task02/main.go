// Задание 2: Вложенный адрес
//
// Тебе нужно:
// 1. Создать структуру Address с полями City (string), Street (string)
// 2. Создать структуру Employee с полями Name (string), Address Address
// 3. В main создать сотрудника и вывести результат
//
// Ожидаемый вывод:
//   Mira: Kazan, Baumana
//
// Запусти: go run main.go

package main


import "fmt"

// TODO: объяви структуру Address с полями City и Street
type Address struct {
	City string
	Street string
}


// TODO: объяви структуру Employee с полями Name и Address
type Employee struct {
	Name string
	Address Address
}

func main() {
    // TODO: создай значение Employee:
    e := Employee{
        Name: "Mira",
        Address: Address{
            City:   "Kazan",
            Street: "Baumana",
        },
    }

    // Name: "Mira", Address: {City: "Kazan", Street: "Baumana"}

    fmt.Printf("%s: %s, %s\n", e.Name, e.Address.City, e.Address.Street)
}