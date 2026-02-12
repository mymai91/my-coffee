package menus

import (
	"context"

	"connectrpc.com/connect"
	menupb "github.com/jany/my-coffee/gen/proto/menu"
	"github.com/jany/my-coffee/gen/proto/menu/menuconnect"
)

// Compile-time check that Server implements the Connect RPC handler interface.
var _ menuconnect.MenuServiceHandler = (*Server)(nil)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) GetMenu(ctx context.Context, req *connect.Request[menupb.GetMenuRequest]) (*connect.Response[menupb.GetMenuResponse], error) {
	items := []*menupb.MenuItem{
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

	return connect.NewResponse(&menupb.GetMenuResponse{
		Items: items,
	}), nil
}