package services

import (
	"github.com/jany/my-coffee/internal/core/domain"
	"github.com/jany/my-coffee/internal/core/ports"
)

// menuService implements ports.MenuService.
type menuService struct{}

// Compile-time check.
var _ ports.MenuService = (*menuService)(nil)

// NewMenuService creates a new MenuService.
func NewMenuService() ports.MenuService {
	return &menuService{}
}

func (s *menuService) GetMenuItems() ([]domain.MenuItem, error) {
	// In a real app this might come from a repository port.
	// For now the menu is hard-coded (same behaviour as original).
	items := []domain.MenuItem{
		{Name: "Espresso", Description: "Strong and rich Italian-style coffee", Price: 2.50},
		{Name: "Latte", Description: "Espresso with steamed milk and a light layer of foam", Price: 3.50},
		{Name: "Cortado", Description: "Equal parts espresso and steamed milk", Price: 3.25},
		{Name: "Ice Latte", Description: "Espresso with cold milk and ice", Price: 3.75},
	}
	return items, nil
}
