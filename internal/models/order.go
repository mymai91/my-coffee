package models

import "time"

type OrderStatus string

const (
	StatusQueued   OrderStatus = "QUEUED"
	StatusGrinding OrderStatus = "GRINDING"
	StatusBrewing  OrderStatus = "BREWING"
	StatusFrothing OrderStatus = "FROTHING"
	StatusReady    OrderStatus = "READY"
)
type Order struct {
	ID uint `gorm:"primaryKey"`
	MenuItemName string `gorm:"not null"`
	Status OrderStatus `gorm:"default:QUEUED"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Order) TableName() string {
	return "orders"
}