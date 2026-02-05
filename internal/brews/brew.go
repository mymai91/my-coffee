package brews

import (
	"context"
	"fmt"
	"log"

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
	log.Printf("OrderDink brew go: %v", req)

	order := &models.Order{
		MenuItemName: req.MenuItemName,
	}

	if err := s.orderRepo.Create(order); err != nil {
		log.Printf("Failed to create order: %v", err)
		return nil, fmt.Errorf("Failed to create order: %w", err)
	}

	return &brewpb.OrderResponse{
		OrderId: fmt.Sprintf("order-%d", order.ID),
	}, nil
}

func (s *Server) ListOrders(ctx context.Context, req *brewpb.ListOrdersRequest) (*brewpb.ListOrdersResponse, error) {
	orders, err := s.orderRepo.FindAll()

	if err != nil {
		log.Printf("Failed to list orders: %v", err)
		return nil, fmt.Errorf("Failed to list orders: %w", err)
	}

	var orderpbs []*brewpb.Order

	for _, order := range orders {
		orderpbs = append(orderpbs, &brewpb.Order{
			OrderId: fmt.Sprintf("order-%d", order.ID),
			MenuItemName: order.MenuItemName,
			Status: string(order.Status),
		})
	}

	return &brewpb.ListOrdersResponse{
		Orders: orderpbs,
	}, nil
}