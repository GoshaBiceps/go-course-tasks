// Задача: Базовый интерфейс
//
// Ожидаемый вывод:
//   Hello, Anna
//   Welcome, stranger

package main

import "fmt"

// TODO: объяви интерфейс Greeter с методом Greet() string
type Greeter interface{
	Greet() string
}

// TODO: объяви структуру User с полем Name string
type User struct {
	Name string 
}
// TODO: реализуй метод Greet() string для User — возвращай "Hello, " + u.Name
func(u User) Greet() string{
	return "Hello," + u.Name
}

// TODO: объяви структуру Guest (без полей)
type Guest struct {

}
// TODO: реализуй метод Greet() string для Guest — возвращай "Welcome, stranger"
func( g Guest) Greet() string {
	return "Welcome, strange"
}
// TODO: напиши функцию printGreeting(g Greeter), которая печатает g.Greet()
func printGreeting(g Greeter) {
	fmt.Println(g.Greet())
}

func main() {
	// TODO: вызови printGreeting для User{Name: "Anna"} и Guest{}
	printGreeting(User{Name: "Anna"})
	printGreeting(Guest{})

	fmt.Println("TODO: implement me")
}
