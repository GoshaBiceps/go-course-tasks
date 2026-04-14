// Задача 6: Интерфейс для обработки платежа
//
// Ожидаемый вывод:
//   paid 100 by card
//   paid 50 in cash

package main

import "fmt"

// TODO: объяви интерфейс PaymentProcessor с методом Pay(amount int) error
type PaymentProcessor interface {
	Pay(amount int) error
}

// TODO: объяви структуру CardProcessor (без полей)
type CardProcessor struct {}

// TODO: реализуй Pay(amount int) error для CardProcessor
//       выводи "paid <amount> by card" и возвращай nil
func (c CardProcessor) Pay(amount int) error {
	fmt.Printf("paid %d by card\n", amount)
	return nil 
}

// TODO: объяви структуру CashProcessor (без полей)
type CashProcessor struct {}

// TODO: реализуй Pay(amount int) error для CashProcessor
//       выводи "paid <amount> in cash" и возвращай nil
func (h CashProcessor) Pay(amount int) error {
	fmt.Printf("paid %d in cash\n", amount)
	return nil 
}

// TODO: напиши функцию checkout(p PaymentProcessor, amount int)
//       вызывает p.Pay(amount), при ошибке выводит её
func checkout(p PaymentProcessor, amount int) {
	err:= p.Pay(amount)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	// TODO: вызови checkout с CardProcessor{} и amount 100
	checkout(CardProcessor{} , 100)

	// TODO: вызови checkout с CashProcessor{} и amount 50
	checkout(CashProcessor{}, 50)
}
