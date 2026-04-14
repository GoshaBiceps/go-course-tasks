// Задача 3: Сравнение двух подходов
//
// Ожидаемый вывод:
//   after copy deposit: 0
//   after pointer deposit: 100

package main

import "fmt"

// TODO: объяви структуру Wallet с полем Balance (int)
type Wallet struct {
	Balance int 
}
// TODO: добавь метод DepositCopy(amount int) с получателем-значением
//       (не меняет исходный объект)
func (w Wallet) DepositCopy(amount int) {
	w.Balance += amount
}
// TODO: добавь метод Deposit(amount int) с получателем-указателем
//       (увеличивает Balance на amount)
func (w *Wallet) Deposit(amount int) {
	w.Balance += amount
}

func main() {
	w := Wallet{Balance: 0}

	w.DepositCopy(100)
	fmt.Println("after copy deposit:", w.Balance)

	w.Deposit(100)
	fmt.Println("after pointer deposit:", w.Balance)
}
