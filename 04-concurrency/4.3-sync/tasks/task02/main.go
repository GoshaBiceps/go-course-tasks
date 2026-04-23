// Задание 2: Параллельные задачи с errgroup
//
// У тебя есть список "сервисов", к которым нужно обратиться параллельно.
// Один из них - намеренно неработающий (возвращает ошибку).
//
// Используй errgroup для параллельного выполнения.
// Если хотя бы один сервис вернул ошибку - выведи её и заверши.
//
// Установи зависимость:
//   go get golang.org/x/sync
//
// Ожидаемый вывод:
//   Опрашиваем: users
//   Опрашиваем: products
//   Опрашиваем: broken-service
//   Ошибка: сервис 'broken-service' недоступен
//
// Порядок строк "Опрашиваем:" может быть любым.
//
// Запусти: go run main.go

package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// callService имитирует обращение к сервису.
// Для "broken-service" возвращает ошибку.
func callService(name string) error {
	fmt.Println("Опрашиваем:", name)
	time.Sleep(100 * time.Millisecond) // имитируем сетевой запрос

	if name == "broken-service" {
		return fmt.Errorf("сервис '%s' недоступен", name)
	}
	return nil
}

func main() {
	services := []string{"users", "products", "broken-service"}

	// TODO: создай errgroup
	var g errgroup.Group

	// TODO: для каждого сервиса запусти g.Go(func() error { ... })
	// Внутри вызови callService(service) и верни результат
	for _, s := range services {
		service := s
		g.Go(func() error { return callService(service) })
	}

	// TODO: вызови g.Wait() и обработай ошибку
	err := g.Wait()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
}
