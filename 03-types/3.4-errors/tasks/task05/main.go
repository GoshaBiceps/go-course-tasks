// Задача 5: Кастомная ошибка
//
// Ожидаемый вывод:
//   validation error: field email: empty value

package main

import (
	"errors"
	"fmt"
)

// TODO: объяви тип InputError со структурой:
//       Field  string
//       Reason string
type InputError struct {
	Field string
	Reason string 
}
// TODO: реализуй метод Error() string
//       возвращай "validation error: field <Field>: <Reason>"
func ( i InputError) Error() string {
	return fmt.Sprintf("validation error: field %s: %s", i.Field, i.Reason)
}
// TODO: напиши функцию validateEmail(email string) error
//       если email == "", верни InputError{Field: "email", Reason: "empty value"}
//       иначе верни nil
func validateEmail(email string) error {
	if email == "" {
		return InputError{Field: "email", Reason: "empty value"}
	}

	return nil 
}
func main() {
	err := validateEmail("")
	if err != nil {
		var inputErr InputError
		if errors.As(err, &inputErr) {
			fmt.Println(inputErr)
		}
	}
}
