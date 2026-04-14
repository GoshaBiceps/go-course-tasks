// Задача: Sentinel error
//
// Ожидаемый вывод:
//   buy ok
//   out of stock: show restock page

package main

import(
"fmt"
"errors"
)

// TODO: объяви sentinel error:
//       var ErrOutOfStock = errors.New("out of stock")
//       (не забудь импортировать "errors")
 var ErrOutOfStock = errors.New("out of stock")

// TODO: напиши функцию buyItem(count int) error
//       если count == 0, верни ErrOutOfStock
//       иначе верни nil
func buyItem(count int) error{
	if count == 0{
		return ErrOutOfStock
	}

	return nil 
}
func main() {
	// TODO: вызови buyItem(5)
	//       если ошибки нет — выведи "buy ok"
     err := buyItem(5)
	 if err == nil{
		 fmt.Println("buy ok")
	 }

	// TODO: вызови buyItem(0)
	//       если errors.Is(err, ErrOutOfStock) — выведи "out of stock: show restock page"
	 err = buyItem(0)
	 if errors.Is(err,ErrOutOfStock) {
			fmt.Println("out of stock: show restock page")
		}
}
