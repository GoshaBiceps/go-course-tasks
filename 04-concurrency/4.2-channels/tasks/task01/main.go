// Задание 1: Worker Pool
//
// Реализуй пул из 3 воркеров, которые обрабатывают 12 задач.
// Задача - число (от 1 до 12). Воркер возводит его в квадрат.
//
// Схема:
//   [задачи 1..12] --> [канал jobs] --> [3 воркера параллельно] --> [канал results] --> [вывод]
//
// Шаги:
//   1. Создай канал jobs (буферизованный, размер 12) и results (буферизованный, размер 12)
//   2. Запусти 3 горутины-воркера через цикл
//      Каждый воркер читает из jobs через range и пишет результат в results
//   3. Закинь задачи 1..12 в jobs, потом вызови close(jobs)
//   4. Дождись завершения воркеров (WaitGroup), потом закрой results
//   5. Выведи все результаты из results
//
// Ожидаемый вывод (порядок любой):
//   Задача 1: 1^2 = 1
//   Задача 3: 3^2 = 9
//   ...
//   Все результаты получены!
//
// Запусти: go run main.go

package main

import (
	"fmt"
	"sync"
)

// TODO: напиши функцию worker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup)
func worker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		square := job * job
		res := fmt.Sprintf("Задача %d: %d^2 = %d", job, job, square)
		results <- res
	}

}

func main() {
	wg := &sync.WaitGroup{}
	// TODO: создай каналы jobs и results
	jobs := make(chan int, 12)
	results := make(chan string)
	// TODO: запусти 3 воркера через цикл
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(i, jobs, results, wg)
	}
	// TODO: закинь задачи 1..12 в jobs и закрой канал
	for i := 1; i <= 12; i++ {
		jobs <- i
	}
	close(jobs)
	// TODO: запусти горутину, которая ждёт все wg.Wait() и потом закрывает results
	go func() {
		wg.Wait()
		close(results)
	}()
	// TODO: выведи результаты из results через range
	for res := range results {
		fmt.Println(res)
	}

	fmt.Println("Все результаты получены!")
}
