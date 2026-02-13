package main

import (
	"log"
	"net/http"

	"github.com/jany/my-coffee/gen/proto/menu/menuconnect"
	"github.com/jany/my-coffee/internal/adapters/handler"
	"github.com/jany/my-coffee/internal/core/services"
)

// cors middleware to allow requests from the Vite dev server
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Wire hexagonal layers: service (core) → handler (driving adapter)
	menuSvc := services.NewMenuService()
	menuHandler := handler.NewMenuHandler(menuSvc)

	mux := http.NewServeMux()
	path, h := menuconnect.NewMenuServiceHandler(menuHandler)
	mux.Handle(path, h)

	// Use h2c so we can serve HTTP/2 without TLS (needed for gRPC compatibility)
	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	server := http.Server{
		Addr:      ":50052",
		Handler:   cors(mux),
		Protocols: p,
	}

	log.Println("MenuService (Connect RPC) running on :50052")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}