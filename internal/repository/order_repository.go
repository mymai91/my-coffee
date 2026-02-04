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

func (r *OrderRepository) FindByID(id uint) (*models.Order, error) {
    var order models.Order
    err := r.db.First(&order, id).Error
    return &order, err
}

func (r *OrderRepository) UpdateStatus(id uint, status models.OrderStatus) error {
    return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}