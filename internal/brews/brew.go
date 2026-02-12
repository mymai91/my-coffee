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
		fmt.Printf("Failed to create order: %v\n", err)
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

func (s *Server) GetOrder(ctx context.Context, req *brewpb.GetOrderRequest) (*brewpb.GetOrderResponse, error) {
	var orderID uint
	if _, err := fmt.Sscanf(req.OrderId, "order-%d", &orderID); err != nil {
		return nil, fmt.Errorf("invalid order ID format: %w", err)
	}

	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &brewpb.GetOrderResponse{
		Order: &brewpb.Order{
			OrderId:      fmt.Sprintf("order-%d", order.ID),
			MenuItemName: order.MenuItemName,
			Status:       string(order.Status),
		},
	}, nil
}

func (s *Server) UpdateOrderStatus(ctx context.Context, req *brewpb.UpdateOrderStatusRequest) (*brewpb.UpdateOrderStatusResponse, error) {
	var orderID uint
	if _, err := fmt.Sscanf(req.OrderId, "order-%d", &orderID); err != nil {
		return nil, fmt.Errorf("invalid order ID format: %w", err)
	}

	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		log.Printf("Failed to find order: %v", err)
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	// Convert proto status to model status
	order.Status = models.OrderStatus(req.Status.String())

	if err := s.orderRepo.Update(order); err != nil {
		log.Printf("Failed to update order status: %v", err)
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return &brewpb.UpdateOrderStatusResponse{
		Order: &brewpb.Order{
			OrderId:      fmt.Sprintf("order-%d", order.ID),
			MenuItemName: order.MenuItemName,
			Status:       string(order.Status),
		},
	}, nil
}

func (s *Server) DeleteOrder(ctx context.Context, req *brewpb.DeleteOrderRequest) (*brewpb.DeleteOrderResponse, error) {
	var orderID uint
	if _, err := fmt.Sscanf(req.OrderId, "order-%d", &orderID); err != nil {
		return nil, fmt.Errorf("invalid order ID format: %w", err)
	}

	if err := s.orderRepo.Delete(orderID); err != nil {
		log.Printf("Failed to delete order: %v", err)
		return nil, fmt.Errorf("failed to delete order: %w", err)
	}

	return &brewpb.DeleteOrderResponse{
		Success: true,
	}, nil
}