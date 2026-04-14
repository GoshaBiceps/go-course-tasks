// Задача 7: Интеграционная задача — система логирования
//
// Ожидаемый вывод:
//   processing order order-1
//   validating order
//   order done
//   [INFO] processing order order-1
//   [INFO] validating order
//   [INFO] order done

package main

import "fmt"

// TODO: объяви интерфейс Logger с методом Log(message string)
type Logger interface {
	Log(message string)
}

// TODO: объяви структуру ConsoleLogger (без полей)
type ConsoleLogger struct {}

// TODO: реализуй Log(message string) — просто выводи message
func ( c ConsoleLogger) Log(message string) {
	fmt.Println(message)
}
// TODO: объяви структуру PrefixLogger с полем Prefix (string)
type PrefixLogger struct {
	Prefix string
}
// TODO: реализуй Log(message string) — выводи "<Prefix> <message>"
func (p PrefixLogger) Log(message string) {
	fmt.Printf("%s %s\n", p.Prefix, message)
}
// TODO: напиши функцию processOrder(logger Logger, id string)
//       пишет 3 сообщения:
//       "processing order <id>"
//       "validating order"
//       "order done"
func processOrder(logger Logger, id string){
	logger.Log(fmt.Sprintf("processing order %s", id))
	logger.Log("validating order")
	logger.Log("order done")
}

func main() {
	// TODO: вызови processOrder с ConsoleLogger{} и id "order-1"
	processOrder(ConsoleLogger{}, "order-1")

	// TODO: вызови processOrder с PrefixLogger{Prefix: "[INFO]"} и id "order-1"
	processOrder(PrefixLogger{Prefix: "[INFO]"}, "order-1")
}
