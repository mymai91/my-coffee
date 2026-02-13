package ports

import "github.com/jany/my-coffee/internal/core/domain"

// OrderService is a DRIVING port — it describes what the application
// can do, and is called by adapters like gRPC handlers or CLI tools.
type OrderService interface {
	CreateOrder(menuItemName string) (*domain.Order, error)
	ListOrders() ([]domain.Order, error)
	GetOrder(id uint) (*domain.Order, error)
	UpdateOrderStatus(id uint, status domain.OrderStatus) (*domain.Order, error)
	DeleteOrder(id uint) error
}

// MenuService is a DRIVING port for menu operations.
type MenuService interface {
	GetMenuItems() ([]domain.MenuItem, error)
}
