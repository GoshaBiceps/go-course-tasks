// Задача 7: Интеграционная задача — регистрация пользователя
//
// Ожидаемый вывод:
//   registered: Alice (25)
//   error: register: name is empty
//   error: register: age must be >= 18

package main

import (
	"fmt"
	"errors"
)

// TODO: напиши функцию validateName(name string) error
//       если name == "", верни ошибку "name is empty"
func validateName(name string) error {
	if name == ""{
		return errors.New("name is empty")
	}

	return nil 
}
// TODO: напиши функцию validateAge(age int) error
//       если age < 18, верни ошибку "age must be >= 18"
func validateAge(age int) error {
	if age < 18 {
		return errors.New("age must be >= 18")
	} 

	return nil 
}
// TODO: напиши функцию register(name string, age int) error
//       вызывает validateName и validateAge
//       оборачивает каждую ошибку через fmt.Errorf: "register: %w"
//       если всё ок — выводи "registered: <name> (<age>)" и возвращай nil
func register(name string , age int) error {
	err := validateName(name)
	if err != nil {
		return fmt.Errorf("registr: %w", err)
	}

	err = validateAge(age)
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	fmt.Printf("registered: %s (%d)\n", name, age)
	return nil 
}

func main() {
	scenarios := []struct {
		name string
		age  int
	}{
		{"Alice", 25},
		{"", 25},
		{"Bob", 15},
	}

	for _, s := range scenarios {
		if err := register(s.name, s.age); err != nil {
			fmt.Println("error:", err)
		}
	}
}
