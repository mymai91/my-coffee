package main

import (
	"log"
	"net/http"

	"github.com/jany/my-coffee/gen/proto/menu/menuconnect"
	"github.com/jany/my-coffee/internal/menus"
)

func main() {
	mux := http.NewServeMux()
	path, handler := menuconnect.NewMenuServiceHandler(menus.New())
	mux.Handle(path, handler)

	// Use h2c so we can serve HTTP/2 without TLS (needed for gRPC compatibility)
	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	server := http.Server{
		Addr:      ":50052",
		Handler:   mux,
		Protocols: p,
	}

	log.Println("MenuService (Connect RPC) running on :50052")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}