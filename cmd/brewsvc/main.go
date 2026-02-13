package main

import (
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/jany/my-coffee/config"
	"github.com/jany/my-coffee/gen/proto/brew/brewconnect"
	"github.com/jany/my-coffee/internal/adapters/handler"
	"github.com/jany/my-coffee/internal/adapters/repository"
	"github.com/jany/my-coffee/internal/core/services"
	database "github.com/jany/my-coffee/internal/datbase"
)

// cors middleware to allow requests from the Vite dev server
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load config
	config.Load()

	// Connect to database
	db := database.Connect()
	defer database.Close()

	// Wire hexagonal layers: repository (driven adapter) → service (core) → handler (driving adapter)
	orderRepo := repository.NewGormOrderRepository(db)
	orderSvc := services.NewOrderService(orderRepo)
	brewHandler := handler.NewBrewHandler(orderSvc)

	// Create Connect RPC server with protovalidate interceptor
	mux := http.NewServeMux()
	path, h := brewconnect.NewBrewServiceHandler(
		brewHandler,
		connect.WithInterceptors(validate.NewInterceptor()),
	)
	mux.Handle(path, h)

	// Use h2c so we can serve HTTP/2 without TLS (needed for gRPC compatibility)
	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	server := http.Server{
		Addr:      ":50051",
		Handler:   cors(mux),
		Protocols: p,
	}

	log.Println("BrewService (Connect RPC) running on :50051")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
