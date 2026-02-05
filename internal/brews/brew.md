
```go
package brews

import (
	"context"
	"fmt"

	brewpb "github.com/jany/my-coffee/gen/proto/brew"
	"github.com/jany/my-coffee/internal/models"
	"github.com/jany/my-coffee/internal/repository"
	"gorm.io/gorm"
)

type Server struct {
	brewpb.UnimplementedBrewServiceServer
	orderRepo *repository.OrderRepository
}

func New(db *gorm.DB) *Server {
	return &Server{
		orderRepo: repository.NewOrderRepository(db),
	}
}

func (s *Server) OrderDrink(ctx context.Context, req *brewpb.OrderRequest) (*brewpb.OrderResponse, error) {
	order := &models.Order{
		MenuItemName: req.MenuItemName,
		Status:        models.StatusQueued,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &brewpb.OrderResponse{
		OrderId: fmt.Sprintf("%d", order.ID),
	}, nil
}

func (s *Server) ListOrders(ctx context.Context, req *brewpb.ListOrdersRequest) (*brewpb.ListOrdersResponse, error) {
	orders, err := s.orderRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	var pbOrders []*brewpb.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, &brewpb.Order{
			OrderId:      fmt.Sprintf("%d", order.ID),
			MenuItemName: order.MenuItemName,
			Status:       mapStatusToProto(order.Status),
		})
	}

	return &brewpb.ListOrdersResponse{
		Orders: pbOrders,
	}, nil
}

func mapStatusToProto(status models.OrderStatus) brewpb.DrinkStatus {
	switch status {
	case models.StatusQueued:
		return brewpb.DrinkStatus_QUEUED
	case models.StatusGrinding:
		return brewpb.DrinkStatus_GRINDING
	case models.StatusBrewing:
		return brewpb.DrinkStatus_BREWING
	case models.StatusFrothing:
		return brewpb.DrinkStatus_FROTHING
	case models.StatusReady:
		return brewpb.DrinkStatus_READY
	default:
		return brewpb.DrinkStatus_DRINK_STATUS_UNSPECIFIED
	}
}
```