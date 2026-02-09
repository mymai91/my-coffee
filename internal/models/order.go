package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Shared validator instance (like Rails ActiveRecord validations)
var validate = validator.New()

type OrderStatus string

const (
	StatusQueued   OrderStatus = "QUEUED"
	StatusGrinding OrderStatus = "GRINDING"
	StatusBrewing  OrderStatus = "BREWING"
	StatusFrothing OrderStatus = "FROTHING"
	StatusReady    OrderStatus = "READY"
)

type Order struct {
	ID           uint        `gorm:"primaryKey"`
	MenuItemName string      `gorm:"not null" validate:"required,min=1"`  // Like: validates :menu_item_name, presence: true, length: { minimum: 1 }
	Status       OrderStatus `gorm:"default:QUEUED" validate:"omitempty,oneof=QUEUED GRINDING BREWING FROTHING READY"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Order) TableName() string {
	return "orders"
}

// BeforeCreate hook - runs validation before saving (like Rails before_create callback)
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	return o.Validate()
}

// BeforeUpdate hook - runs validation before updating (like Rails before_update callback)
func (o *Order) BeforeUpdate(tx *gorm.DB) error {
	return o.Validate()
}

// Validate runs all validations (like Rails valid? method)
func (o *Order) Validate() error {
	return validate.Struct(o)
}