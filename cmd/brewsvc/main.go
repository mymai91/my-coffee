package main

import (
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/jany/my-coffee/config"
	"github.com/jany/my-coffee/gen/proto/brew/brewconnect"
	"github.com/jany/my-coffee/internal/brews"
	database "github.com/jany/my-coffee/internal/datbase"
)

func main() {
	// Load config
	config.Load()

	// Connect to database
	db := database.Connect()
	defer database.Close()

	// Create Connect RPC server with protovalidate interceptor
	mux := http.NewServeMux()
	path, handler := brewconnect.NewBrewServiceHandler(
		brews.New(db),
		connect.WithInterceptors(validate.NewInterceptor()),
	)
	mux.Handle(path, handler)

	// Use h2c so we can serve HTTP/2 without TLS (needed for gRPC compatibility)
	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	server := http.Server{
		Addr:      ":50051",
		Handler:   mux,
		Protocols: p,
	}

	log.Println("BrewService (Connect RPC) running on :50051")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
