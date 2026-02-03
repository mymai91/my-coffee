package menus

import (
	"context"

	menupb "github.com/jany/my-coffee/gen/proto/menu"
)


type Server struct {
	menupb.UnimplementedMenuServiceServer
}

func New() *Server {
	return &Server{}
}

func (s *Server) GetMenu(ctx context.Context, req *menupb.GetMenuRequest) (*menupb.GetMenuResponse, error) {
	item := []*menupb.MenuItem{
		{
			Name:        "Espresso",
			Description: "Strong and rich Italian-style coffee",
			Price:       2.50,
		},
		{
			Name:        "Latte",
			Description: "Espresso with steamed milk and a light layer of foam",
			Price:       3.50,
		},
		{
			Name:        "Cortado",
			Description: "Equal parts espresso and steamed milk",
			Price:       3.25,
		},
		{
			Name:        "Ice Latte",
			Description: "Espresso with cold milk and ice",
			Price:       3.75,
		},
	}

	return &menupb.GetMenuResponse{
		Items: item,
	}, nil
}