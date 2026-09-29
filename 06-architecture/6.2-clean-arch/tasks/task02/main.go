// Задание 1: Repository pattern для товаров
//
// В этой папке три пакета: domain, repository, service.
// В каждом из них есть файл с TODO-комментариями.
// Реализуй их по порядку: domain -> repository -> service.
//
// После реализации запусти main.go - он проверит что всё работает.
//
// Ожидаемый вывод:
//   Создан: iPhone 15 (цена: 99999.00, на складе: 10)
//   Создан: AirPods (цена: 14999.00, на складе: 5)
//   Покупаем 3 iPhone 15... успешно
//   На складе iPhone 15 осталось: 7
//   Пробуем купить 10 AirPods: недостаточно товара на складе
//
// Запусти: go run .

package main

import (
	"fmt"
	"log"

	"github.com/go-course/clean-arch-task01/repository"
	"github.com/go-course/clean-arch-task01/service"
	"go.uber.org/fx"
)

func run(svc *service.ProductService) {
	p1, err := svc.Create("iPhone 15", 99999.00, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"Создан: %s (цена: %.2f, на складе: %d)\n",
		p1.Name, p1.Price, p1.Stock,
	)

	p2, err := svc.Create("AirPods", 14999.00, 5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"Создан: %s (цена: %.2f, на складе: %d)\n",
		p2.Name, p2.Price, p2.Stock,
	)

	fmt.Println("Покупаем 3 iPhone 15...")
	if err := svc.Buy(p1.ID, 3); err != nil {
		fmt.Println("ошибка:", err)
	} else {
		fmt.Println("успешно")
	}

	updated, _ := svc.List()

	for _, p := range updated {
		if p.ID == p1.ID {
			fmt.Printf(
				"На складе %s осталось: %d\n",
				p.Name, p.Stock,
			)
		}
	}

	fmt.Println("Пробуем купить 10 AirPods:")
	if err := svc.Buy(p2.ID, 10); err != nil {
		fmt.Println(err)
	}
}

func main() {
	app := fx.New(
		fx.Provide(
			fx.Annotate(
				repository.NewInMemoryProductRepository,
				fx.As(new(repository.ProductRepository)),
			),
			service.NewProductService,
		),
		fx.Invoke(run),
	)

	app.Run()
}
