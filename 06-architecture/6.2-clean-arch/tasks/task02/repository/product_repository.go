package repository

import (
	"errors"
	"sync"

	"github.com/go-course/clean-arch-task01/domain"
)

// TODO: создай интерфейс ProductRepository с методами:
//
//	Save(product domain.Product) (domain.Product, error)
//	FindAll() ([]domain.Product, error)
//	FindByID(id int) (domain.Product, error)
//	Delete(id int) error
//	UpdateStock(id int, newStock int) error
type ProductRepository interface {
	Save(product domain.Product) (domain.Product, error)
	FindAll() ([]domain.Product, error)
	FindByID(id int) (domain.Product, error)
	Delete(id int) error
	UpdateStock(id int, newStock int) error
}

// TODO: создай структуру InMemoryProductRepository и реализуй все методы интерфейса
// Используй map[int]domain.Product для хранения и sync.RWMutex для защиты от гонок
type InMemoryProductRepository struct {
	products map[int]domain.Product
	mu       sync.RWMutex
	nextID   int
}

// TODO: создай конструктор NewInMemoryProductRepository() *InMemoryProductRepository
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[int]domain.Product),
	}
}

func (r *InMemoryProductRepository) Save(product domain.Product) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	product.ID = r.nextID

	r.products[product.ID] = product // сохранили под наш айдишник наш продукт

	return product, nil
}

func (r *InMemoryProductRepository) FindByID(id int) (domain.Product, error) { // поиск продукта через айди
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[id]
	if !ok {
		return domain.Product{}, errors.New(" Product not found!")
	}

	return product, nil
}

func (r *InMemoryProductRepository) FindAll() ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]domain.Product, 0, len(r.products))

	for _, product := range r.products {
		products = append(products, product)
	}
	return products, nil
}

func (r *InMemoryProductRepository) Delete(id int) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.products[id]
	if !ok {
		return errors.New("product not faund")
	}

	delete(r.products, id) // удаляем продукт из мапы

	return nil
}

func (r *InMemoryProductRepository) UpdateStock(id int, newStock int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	product, ok := r.products[id]
	if !ok {
		return errors.New("product not found")
	}

	product.Stock = newStock // обновилди сток
	r.products[id] = product // записали в мапу

	return nil
}
