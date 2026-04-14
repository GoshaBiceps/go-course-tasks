// Задача 3: Проверка реализации через компилятор
//
// Ожидаемый вывод:
//   Running fast

package main

import "fmt"

// TODO: объяви интерфейс Runner с методом Run() string
type Runner interface {
	Run() string
}
// TODO: объяви структуру Athlete (без полей)
type Athlete struct{

}
// TODO: реализуй метод Run() string для Athlete — возвращай "Running fast"
func (a Athlete) Run() string {
	return "Running fast"
}

// Проверка через компилятор: раскомментируй строку ниже после реализации

var _ Runner = Athlete{}

func main() {
	// TODO: создай Athlete и вызови Run() через интерфейс Runner
	a := Athlete{}

	var r Runner = a 
	
	fmt.Println(r.Run())
	
}