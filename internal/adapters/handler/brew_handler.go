package handler

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	brewpb "github.com/jany/my-coffee/gen/proto/brew"
	"github.com/jany/my-coffee/gen/proto/brew/brewconnect"
	"github.com/jany/my-coffee/internal/core/domain"
	"github.com/jany/my-coffee/internal/core/ports"
)

// BrewHandler is a driving adapter that translates Connect RPC requests
// into calls on the OrderService port.
type BrewHandler struct {
	svc ports.OrderService
}

// Compile-time check that BrewHandler implements the Connect RPC handler interface.
var _ brewconnect.BrewServiceHandler = (*BrewHandler)(nil)

// NewBrewHandler creates a new Connect RPC handler for the brew service.
func NewBrewHandler(svc ports.OrderService) *BrewHandler {
	return &BrewHandler{svc: svc}
}

func (h *BrewHandler) OrderDrink(ctx context.Context, req *connect.Request[brewpb.OrderRequest]) (*connect.Response[brewpb.OrderResponse], error) {
	log.Printf("OrderDrink: %v", req.Msg)

	order, err := h.svc.CreateOrder(req.Msg.MenuItemName)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&brewpb.OrderResponse{
		OrderId: fmt.Sprintf("order-%d", order.ID),
	}), nil
}

func (h *BrewHandler) ListOrders(ctx context.Context, req *connect.Request[brewpb.ListOrdersRequest]) (*connect.Response[brewpb.ListOrdersResponse], error) {
	orders, err := h.svc.ListOrders()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var orderpbs []*brewpb.Order
	for _, o := range orders {
		orderpbs = append(orderpbs, toProtoOrder(&o))
	}

	return connect.NewResponse(&brewpb.ListOrdersResponse{
		Orders: orderpbs,
	}), nil
}

func (h *BrewHandler) GetOrder(ctx context.Context, req *connect.Request[brewpb.GetOrderRequest]) (*connect.Response[brewpb.GetOrderResponse], error) {
	orderID, err := parseOrderID(req.Msg.OrderId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	order, err := h.svc.GetOrder(orderID)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&brewpb.GetOrderResponse{
		Order: toProtoOrder(order),
	}), nil
}

func (h *BrewHandler) UpdateOrderStatus(ctx context.Context, req *connect.Request[brewpb.UpdateOrderStatusRequest]) (*connect.Response[brewpb.UpdateOrderStatusResponse], error) {
	orderID, err := parseOrderID(req.Msg.OrderId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	status := domain.OrderStatus(req.Msg.Status.String())
	order, err := h.svc.UpdateOrderStatus(orderID, status)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&brewpb.UpdateOrderStatusResponse{
		Order: toProtoOrder(order),
	}), nil
}

func (h *BrewHandler) DeleteOrder(ctx context.Context, req *connect.Request[brewpb.DeleteOrderRequest]) (*connect.Response[brewpb.DeleteOrderResponse], error) {
	orderID, err := parseOrderID(req.Msg.OrderId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := h.svc.DeleteOrder(orderID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&brewpb.DeleteOrderResponse{
		Success: true,
	}), nil
}

// --- helpers ---

func parseOrderID(raw string) (uint, error) {
	var id uint
	if _, err := fmt.Sscanf(raw, "order-%d", &id); err != nil {
		return 0, fmt.Errorf("invalid order ID format: %w", err)
	}
	return id, nil
}

func toProtoOrder(o *domain.Order) *brewpb.Order {
	return &brewpb.Order{
		OrderId:      fmt.Sprintf("order-%d", o.ID),
		MenuItemName: o.MenuItemName,
		Status:       string(o.Status),
	}
}
