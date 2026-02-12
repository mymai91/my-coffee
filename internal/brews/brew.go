package brews

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	brewpb "github.com/jany/my-coffee/gen/proto/brew"
	"github.com/jany/my-coffee/gen/proto/brew/brewconnect"
	"github.com/jany/my-coffee/internal/models"
	"github.com/jany/my-coffee/internal/repository"
	"gorm.io/gorm"
)

// Compile-time check that Server implements the Connect RPC handler interface.
var _ brewconnect.BrewServiceHandler = (*Server)(nil)

type Server struct {
	orderRepo *repository.OrderRepository
}

func New(db *gorm.DB) *Server {
	return &Server{
		orderRepo: repository.NewOrderRepository(db),
	}
}

func (s *Server) OrderDrink(ctx context.Context, req *connect.Request[brewpb.OrderRequest]) (*connect.Response[brewpb.OrderResponse], error) {
	log.Printf("OrderDrink brew go: %v", req.Msg)

	order := &models.Order{
		MenuItemName: req.Msg.MenuItemName,
	}

	if err := s.orderRepo.Create(order); err != nil {
		log.Printf("Failed to create order: %v", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create order: %w", err))
	}

	return connect.NewResponse(&brewpb.OrderResponse{
		OrderId: fmt.Sprintf("order-%d", order.ID),
	}), nil
}

func (s *Server) ListOrders(ctx context.Context, req *connect.Request[brewpb.ListOrdersRequest]) (*connect.Response[brewpb.ListOrdersResponse], error) {
	orders, err := s.orderRepo.FindAll()
	if err != nil {
		log.Printf("Failed to list orders: %v", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to list orders: %w", err))
	}

	var orderpbs []*brewpb.Order
	for _, order := range orders {
		orderpbs = append(orderpbs, &brewpb.Order{
			OrderId:      fmt.Sprintf("order-%d", order.ID),
			MenuItemName: order.MenuItemName,
			Status:       string(order.Status),
		})
	}

	return connect.NewResponse(&brewpb.ListOrdersResponse{
		Orders: orderpbs,
	}), nil
}

func (s *Server) GetOrder(ctx context.Context, req *connect.Request[brewpb.GetOrderRequest]) (*connect.Response[brewpb.GetOrderResponse], error) {
	var orderID uint
	if _, err := fmt.Sscanf(req.Msg.OrderId, "order-%d", &orderID); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid order ID format: %w", err))
	}

	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("failed to get order: %w", err))
	}

	return connect.NewResponse(&brewpb.GetOrderResponse{
		Order: &brewpb.Order{
			OrderId:      fmt.Sprintf("order-%d", order.ID),
			MenuItemName: order.MenuItemName,
			Status:       string(order.Status),
		},
	}), nil
}

func (s *Server) UpdateOrderStatus(ctx context.Context, req *connect.Request[brewpb.UpdateOrderStatusRequest]) (*connect.Response[brewpb.UpdateOrderStatusResponse], error) {
	var orderID uint
	if _, err := fmt.Sscanf(req.Msg.OrderId, "order-%d", &orderID); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid order ID format: %w", err))
	}

	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		log.Printf("Failed to find order: %v", err)
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("failed to find order: %w", err))
	}

	// Convert proto status to model status
	order.Status = models.OrderStatus(req.Msg.Status.String())

	if err := s.orderRepo.Update(order); err != nil {
		log.Printf("Failed to update order status: %v", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update order status: %w", err))
	}

	return connect.NewResponse(&brewpb.UpdateOrderStatusResponse{
		Order: &brewpb.Order{
			OrderId:      fmt.Sprintf("order-%d", order.ID),
			MenuItemName: order.MenuItemName,
			Status:       string(order.Status),
		},
	}), nil
}

func (s *Server) DeleteOrder(ctx context.Context, req *connect.Request[brewpb.DeleteOrderRequest]) (*connect.Response[brewpb.DeleteOrderResponse], error) {
	var orderID uint
	if _, err := fmt.Sscanf(req.Msg.OrderId, "order-%d", &orderID); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid order ID format: %w", err))
	}

	if err := s.orderRepo.Delete(orderID); err != nil {
		log.Printf("Failed to delete order: %v", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete order: %w", err))
	}

	return connect.NewResponse(&brewpb.DeleteOrderResponse{
		Success: true,
	}), nil
}