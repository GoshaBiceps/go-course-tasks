// Задание 1: Описание книги (value receiver)
//
// Тебе нужно:
// 1. Создать структуру Book с полями Title (string) и Pages (int)
// 2. Добавить метод Summary() string с получателем-значением,
//    который возвращает строку "<Title> (<Pages> pages)"
// 3. В main создать книгу и вывести результат Summary()
//
// Ожидаемый вывод:
//   Go in Action (350 pages)
//
// Запусти: go run main.go

package main

import "fmt"

// TODO: объяви структуру Book с полями Title и Pages
type Book struct{
	Title string
	Pages int 
}
// TODO: добавь метод Summary() string с получателем-значением
func (b Book) Summary() string{
	return fmt.Sprintf("%s (%d pages)", b.Title, b.Pages)
}
func main() {
	// TODO: создай значение Book с Title: "Go in Action", Pages: 350
	b := Book{
		Title: "Go in Action",
		Pages: 350,
	}

	fmt.Println(b.Summary())
}