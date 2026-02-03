package brews

import (
	"context"

	brewpb "github.com/jany/my-coffee/gen/proto/brew"
)

type Server struct {
	brewpb.UnimplementedBrewServiceServer
}

func New() *Server {
	return &Server{}
}

func (s *Server) OrderDrink(ctx context.Context, req *brewpb.OrderRequest) (*brewpb.OrderResponse, error) {
	return &brewpb.OrderResponse{
		OrderId: "12345",
	}, nil
}

func (s *Server) ListOrders(ctx context.Context, req *brewpb.ListOrdersRequest) (*brewpb.ListOrdersResponse, error) {
	orders := []*brewpb.Order{
		{
			OrderId:   "12345",
			MenuItemName: "Espresso",
			Status: brewpb.DrinkStatus_BREWING,
		},
	}

	return &brewpb.ListOrdersResponse{
		Orders: orders,
	},nil
}