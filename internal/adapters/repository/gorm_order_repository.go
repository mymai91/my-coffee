package repository

import (
	"github.com/jany/my-coffee/internal/core/domain"
	"github.com/jany/my-coffee/internal/core/ports"
	"gorm.io/gorm"
)

// GormOrderRepository is a driven adapter that implements ports.OrderRepository
// using GORM as the persistence mechanism.
type GormOrderRepository struct {
	db *gorm.DB
}

// Compile-time check.
var _ ports.OrderRepository = (*GormOrderRepository)(nil)

// NewGormOrderRepository creates a new GORM-backed order repository.
func NewGormOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) Create(order *domain.Order) error {
	model := toModel(order)
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	// Propagate the generated ID back to the domain entity.
	order.ID = model.ID
	order.CreatedAt = model.CreatedAt
	order.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *GormOrderRepository) FindAll() ([]domain.Order, error) {
	var models []OrderModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	orders := make([]domain.Order, len(models))
	for i, m := range models {
		orders[i] = *m.toDomain()
	}
	return orders, nil
}

func (r *GormOrderRepository) FindByID(id uint) (*domain.Order, error) {
	var model OrderModel
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return model.toDomain(), nil
}

func (r *GormOrderRepository) Update(order *domain.Order) error {
	model := toModel(order)
	return r.db.Save(model).Error
}

func (r *GormOrderRepository) Delete(id uint) error {
	return r.db.Delete(&OrderModel{}, id).Error
}
