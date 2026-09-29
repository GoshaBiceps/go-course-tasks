package service

import (
	"errors"

	"github.com/go-course/clean-arch-task01/domain"
	"github.com/go-course/clean-arch-task01/repository"
)

// TODO: создай структуру ProductService с полем repo типа repository.ProductRepository (интерфейс!)
type ProductService struct {
	repo repository.ProductRepository
}

// TODO: создай конструктор NewProductService(repo repository.ProductRepository) *ProductService
func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// TODO: реализуй методы:
//   Create(name string, price float64, stock int) (domain.Product, error)
//     - создаёт продукт через domain.Product, валидирует, сохраняет
//
//   List() ([]domain.Product, error)
//     - возвращает все продукты
//
//   Buy(productID int, quantity int) error
//     - находит продукт, проверяет что quantity <= Stock
//     - если ок - уменьшает Stock через UpdateStock
//     - если нет - возвращает ошибку "недостаточно товара на складе"

func (s *ProductService) Create(name string, price float64, stock int) (domain.Product, error) {

	product := domain.Product{Name: name, Price: price, Stock: stock}

	if err := product.Validate(); err != nil { // так если не прошли валидацию
		return domain.Product{}, err
	}

	product, err := s.repo.Save(product)
	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func (s *ProductService) List() ([]domain.Product, error) {

	products, err := s.repo.FindAll()
	if err != nil {
		return []domain.Product{}, err
	}
	return products, nil
}

func (s *ProductService) Buy(productID int, quantity int) error {
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return err
	}

	if quantity > product.Stock {
		return errors.New("недостаточно товара на складе")
	}

	newStock := product.Stock - quantity

	if err := s.repo.UpdateStock(productID, newStock); err != nil {
		return err
	}

	return nil
}
