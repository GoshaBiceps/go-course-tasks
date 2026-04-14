// Задача 6: Небольшая модель состояния
//
// Ожидаемый вывод:
//   running
//   stopped

package main

import "fmt"

// TODO: объяви структуру Timer с полями Seconds (int), Running (bool)
type Timer struct {
	Seconds int 
	Running bool 
}
// TODO: добавь метод Start() с получателем-указателем — ставит Running = true
func (t *Timer) Start() {
	t.Running = true
}
// TODO: добавь метод Stop() с получателем-указателем — ставит Running = false
func (t *Timer) Stop() {
	t.Running = false
}
// TODO: добавь метод Status() string с получателем-значением
//       возвращает "running" если Running == true, иначе "stopped"
func (t Timer) Status()string {
	if t.Running {
		return "runnig"
	}
	return "stopped"
}

func main() {
	t := &Timer{}
	t.Start()
	fmt.Println(t.Status())
	t.Stop()
	fmt.Println(t.Status())
}
