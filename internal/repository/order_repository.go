package repository

import (
	"github.com/jany/my-coffee/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Find(&orders).Error
	return orders, err
}