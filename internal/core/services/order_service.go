package services

import (
	"fmt"

	"github.com/jany/my-coffee/internal/core/domain"
	"github.com/jany/my-coffee/internal/core/ports"
)

// orderService implements ports.OrderService.
// It contains the business / use-case logic and depends only on port interfaces.
type orderService struct {
	repo ports.OrderRepository
}

// Compile-time check.
var _ ports.OrderService = (*orderService)(nil)

// NewOrderService creates a new OrderService with the given repository port.
func NewOrderService(repo ports.OrderRepository) ports.OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(menuItemName string) (*domain.Order, error) {
	order := domain.NewOrder(menuItemName)

	if err := s.repo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *orderService) ListOrders() ([]domain.Order, error) {
	orders, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	return orders, nil
}

func (s *orderService) GetOrder(id uint) (*domain.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (s *orderService) UpdateOrderStatus(id uint, status domain.OrderStatus) (*domain.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	order.Status = status

	if err := s.repo.Update(order); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return order, nil
}

func (s *orderService) DeleteOrder(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
