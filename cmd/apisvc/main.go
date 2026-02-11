package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	brewpb "github.com/jany/my-coffee/gen/proto/brew"
	menupb "github.com/jany/my-coffee/gen/proto/menu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// cors middleware to allow requests from the Vite dev server
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Connect to Menu Service (port 50052)
	menuConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Cannot connect to menu service: %v", err)
	}
	defer menuConn.Close()

	// Connect to Brew Service (port 50051)
	brewConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Cannot connect to brew service: %v", err)
	}
	defer brewConn.Close()

	menuClient := menupb.NewMenuServiceClient(menuConn)
	brewClient := brewpb.NewBrewServiceClient(brewConn)

	mux := http.NewServeMux()

	// GET /api/menu — returns the coffee menu
	mux.HandleFunc("GET /api/menu", func(w http.ResponseWriter, r *http.Request) {
		resp, err := menuClient.GetMenu(context.Background(), &menupb.GetMenuRequest{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get menu: %v", err), http.StatusInternalServerError)
			return
		}

		type MenuItem struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
		}

		items := make([]MenuItem, 0, len(resp.Items))
		for _, item := range resp.Items {
			items = append(items, MenuItem{
				Name:        item.Name,
				Description: item.Description,
				Price:       item.Price,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	// GET /api/orders — returns all orders
	mux.HandleFunc("GET /api/orders", func(w http.ResponseWriter, r *http.Request) {
		resp, err := brewClient.ListOrders(context.Background(), &brewpb.ListOrdersRequest{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to list orders: %v", err), http.StatusInternalServerError)
			return
		}

		type Order struct {
			OrderID      string `json:"orderId"`
			MenuItemName string `json:"menuItemName"`
			Status       string `json:"status"`
		}

		orders := make([]Order, 0, len(resp.Orders))
		for _, o := range resp.Orders {
			orders = append(orders, Order{
				OrderID:      o.OrderId,
				MenuItemName: o.MenuItemName,
				Status:       o.Status,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	})

	// POST /api/orders — create a new order
	mux.HandleFunc("POST /api/orders", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			MenuItemName string `json:"menuItemName"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.MenuItemName == "" {
			http.Error(w, "menuItemName is required", http.StatusBadRequest)
			return
		}

		resp, err := brewClient.OrderDrink(context.Background(), &brewpb.OrderRequest{
			MenuItemName: req.MenuItemName,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to order drink: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"orderId": resp.OrderId,
		})
	})

	log.Println("API Service running on :9000")
	if err := http.ListenAndServe(":9000", cors(mux)); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
