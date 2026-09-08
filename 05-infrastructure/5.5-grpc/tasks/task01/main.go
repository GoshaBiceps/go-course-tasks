package main

import (
	"context"
	"errors"
	"fmt"
)

type CreateOrderRequest struct {
	CustomerID string
	Items      []string
}

type CreateOrderResponse struct {
	OrderID string
	Status  string
}

type GetOrderRequest struct {
	OrderID string
}

type GetOrderResponse struct {
	OrderID    string
	CustomerID string
	Items      []string
	Status     string
}

var ErrOrderNotFound = errors.New("order not found")

type OrderServiceServer interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error)
}

// TODO: Implement orderServiceImpl.
// - Store orders in a map keyed by order ID.
// - Generate IDs like "order-1", "order-2", etc.
// - CreateOrder: validate CustomerID and Items are non-empty; store and return status "created".
// - GetOrder: return ErrOrderNotFound for unknown IDs.

type order struct {
	id         string
	customerID string
	items      []string
	status     string
}

type orderServiceImpl struct {
	orders map[string]order
	nextID int
}

func NewOrderService() OrderServiceServer {
	return &orderServiceImpl{
		orders: make(map[string]order),
		nextID: 1,
	}
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	// TODO: validate input, store order, return response
	if req.CustomerID == "" {
		return nil, errors.New("customer_id is required")
	}

	if len(req.Items) == 0 {
		return nil, errors.New("items are required")
	}

	orderID := fmt.Sprintf("order-%d", s.nextID) // создали id нового заказа

	newOrder := order{ //  сформировали сам заказ
		id:         orderID,
		customerID: req.CustomerID,
		items:      req.Items,
		status:     "created",
	}

	s.orders[newOrder.id] = newOrder // положили заказ в хранилище
	s.nextID++                       // увеличиваем счетчик

	orederResponse := &CreateOrderResponse{ // собрали ответ клиенту
		OrderID: newOrder.id,
		Status:  newOrder.status,
	}

	return orederResponse, nil
}

func (s *orderServiceImpl) GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error) {
	// TODO: look up order by ID, return ErrOrderNotFound if missing
	foundOrder, ok := s.orders[req.OrderID] // нашли по айди заказ
	if !ok {
		return nil, ErrOrderNotFound
	}

	orderResponse := &GetOrderResponse{ //  собрали ответ клиенту
		OrderID:    foundOrder.id,
		CustomerID: foundOrder.customerID,
		Items:      foundOrder.items,
		Status:     foundOrder.status,
	}

	return orderResponse, nil
}

func main() {
	svc := NewOrderService()
	ctx := context.Background()

	created, err := svc.CreateOrder(ctx, &CreateOrderRequest{
		CustomerID: "c-1",
		Items:      []string{"item-a", "item-b"},
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Created: %+v\n", created)

	got, err := svc.GetOrder(ctx, &GetOrderRequest{OrderID: created.OrderID})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Got: %+v\n", got)
	}

	_, err = svc.GetOrder(ctx, &GetOrderRequest{OrderID: "unknown"})
	if errors.Is(err, ErrOrderNotFound) {
		fmt.Println("Order 'unknown': not found (expected)")
	}
}
