package ports

import "github.com/jany/my-coffee/internal/core/domain"

// OrderRepository is a DRIVEN port — the core tells the outside world
// how it needs data to be persisted, but never how (no GORM, no SQL).
type OrderRepository interface {
	Create(order *domain.Order) error
	FindAll() ([]domain.Order, error)
	FindByID(id uint) (*domain.Order, error)
	Update(order *domain.Order) error
	Delete(id uint) error
}
