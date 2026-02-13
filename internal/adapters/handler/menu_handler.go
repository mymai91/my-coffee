package handler

import (
	"context"

	"connectrpc.com/connect"
	menupb "github.com/jany/my-coffee/gen/proto/menu"
	"github.com/jany/my-coffee/gen/proto/menu/menuconnect"
	"github.com/jany/my-coffee/internal/core/ports"
)

// MenuHandler is a driving adapter that translates Connect RPC requests
// into calls on the MenuService port.
type MenuHandler struct {
	svc ports.MenuService
}

// Compile-time check.
var _ menuconnect.MenuServiceHandler = (*MenuHandler)(nil)

// NewMenuHandler creates a new Connect RPC handler for the menu service.
func NewMenuHandler(svc ports.MenuService) *MenuHandler {
	return &MenuHandler{svc: svc}
}

func (h *MenuHandler) GetMenu(ctx context.Context, req *connect.Request[menupb.GetMenuRequest]) (*connect.Response[menupb.GetMenuResponse], error) {
	items, err := h.svc.GetMenuItems()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var pbItems []*menupb.MenuItem
	for _, item := range items {
		pbItems = append(pbItems, &menupb.MenuItem{
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
		})
	}

	return connect.NewResponse(&menupb.GetMenuResponse{
		Items: pbItems,
	}), nil
}
