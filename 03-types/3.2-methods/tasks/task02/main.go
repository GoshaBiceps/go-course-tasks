// Задание 2: Пополнение кошелька (pointer receiver)
//
// Тебе нужно:
// 1. Создать структуру Wallet с полем Balance (int)
// 2. Добавить метод Deposit(amount int) с получателем-указателем,
//    который увеличивает баланс на amount
// 3. В main создать кошелёк, сделать два пополнения и вывести баланс
//
// Ожидаемый вывод:
//   Balance: 300
//
// Запусти: go run main.go

package main

import "fmt"

// TODO: объяви структуру Wallet с полем Balance
type Wallet struct {
	Balance int 
}

// TODO: добавь метод Deposit(amount int) с получателем-указателем
func (w *Wallet) Deposit(amount int){
	w.Balance += amount
}
func main() {
	// TODO: создай Wallet с Balance: 0
	w := Wallet{Balance: 0}

	w.Deposit(100)
	w.Deposit(200)
	// TODO: вызови Deposit(100) и Deposit(200)

	fmt.Println("Balance:", w.Balance)
}