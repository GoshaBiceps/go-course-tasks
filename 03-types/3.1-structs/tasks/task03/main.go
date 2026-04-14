package main

import "fmt"

// TODO: объяви структуру Address с полями City и Street
type Address struct {
	City   string
	Street string
}

// TODO: объяви структуру Employee с полями Name и Address
type Employee struct {
	Name    string
	Address Address
}

func main() {
	// TODO: создай значение Employee:
	// Name: "Mira", Address: {City: "Kazan", Street: "Baumana"}
	e := Employee{
		Name: "Mira",
		Address: Address{
			City:   "Kazan",
			Street: "Baumana",
		},
	}

	fmt.Printf("%s: %s, %s\n", e.Name, e.Address.City, e.Address.Street)
}