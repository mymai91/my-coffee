package repository

import (
	"time"

	"github.com/jany/my-coffee/internal/core/domain"
)

// OrderModel is the GORM-specific persistence model.
// It lives in the adapter layer so the domain stays free of ORM tags.
type OrderModel struct {
	ID           uint               `gorm:"primaryKey"`
	MenuItemName string             `gorm:"not null"`
	Status       domain.OrderStatus `gorm:"default:QUEUED"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (OrderModel) TableName() string {
	return "orders"
}

// toDomain converts the persistence model to the domain entity.
func (m *OrderModel) toDomain() *domain.Order {
	return &domain.Order{
		ID:           m.ID,
		MenuItemName: m.MenuItemName,
		Status:       m.Status,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// toModel converts a domain entity to the persistence model.
func toModel(o *domain.Order) *OrderModel {
	return &OrderModel{
		ID:           o.ID,
		MenuItemName: o.MenuItemName,
		Status:       o.Status,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}
