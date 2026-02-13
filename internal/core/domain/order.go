package domain

import "time"

// OrderStatus represents the lifecycle status of a coffee order.
type OrderStatus string

const (
	StatusQueued   OrderStatus = "QUEUED"
	StatusGrinding OrderStatus = "GRINDING"
	StatusBrewing  OrderStatus = "BREWING"
	StatusFrothing OrderStatus = "FROTHING"
	StatusReady    OrderStatus = "READY"
)

// Order is the core domain entity — it has zero infrastructure dependencies.
type Order struct {
	ID           uint
	MenuItemName string
	Status       OrderStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewOrder creates a new order with default QUEUED status.
func NewOrder(menuItemName string) *Order {
	return &Order{
		MenuItemName: menuItemName,
		Status:       StatusQueued,
	}
}
