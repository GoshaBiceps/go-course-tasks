// Задача 4: Композиция интерфейсов
//
// Ожидаемый вывод:
//   engine started
//   engine stopped

package main

import "fmt"

// TODO: объяви интерфейс Starter с методом Start()
type Starter interface {
	Start()
}

// TODO: объяви интерфейс Stopper с методом Stop()
type Stopper interface {
	Stop()
}

// TODO: объяви интерфейс Machine как композицию Starter и Stopper
type Machine interface {
	Starter
	Stopper
}
// TODO: объяви структуру Engine (без полей)
type Engine struct {

}
// TODO: реализуй Start() — выводи "engine started"
func (e Engine) Start(){
	fmt.Println("engine started")
}
// TODO: реализуй Stop()  — выводи "engine stopped"
func (e Engine) Stop() {
	fmt.Println("engine stopped")
}
// TODO: напиши функцию runCycle(m Machine)
//       вызывает m.Start() и m.Stop()
func runCycle(m Machine) {

	m.Start()
	m.Stop()
}

func main() {
	// TODO: вызови runCycle с Engine{}
	e := Engine{}
	runCycle(e)
}
