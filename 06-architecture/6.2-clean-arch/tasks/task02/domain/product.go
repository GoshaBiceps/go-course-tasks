package domain

import "errors"

// TODO: создай тип Product со следующими полями:
//
//	ID    int
//	Name  string
//	Price float64
//	Stock int    ← количество на складе
type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

// TODO: добавь метод Validate() error:
//   - Name не должен быть пустым
//   - Price должна быть больше нуля
//   - Stock должен быть >= 0
func (p Product) Validate() error {

	if p.Name == "" {
		return errors.New("name is empty")
	}

	if p.Price <= 0 {
		return errors.New("price less than zero")
	}

	if p.Stock <= 0 {
		return errors.New("insufficient quantity of goods")
	}

	return nil
}
