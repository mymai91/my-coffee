package main

import (
	"fmt"
	"net"
	"os"

	brewpb "github.com/jany/my-coffee/gen/proto/brew"
	"github.com/jany/my-coffee/internal/brews"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Failed to listen on Port 50051: %v\n", err)
		return
	}

	server := grpc.NewServer()
	brewService := brews.New()
	brewpb.RegisterBrewServiceServer(server, brewService)

	// Only enable reflection in development
    if os.Getenv("ENV") != "production" {
        reflection.Register(server)
    }

	fmt.Println("Menu Service is running on port 50051...")
	if err := server.Serve(listener); err != nil {
		fmt.Printf("Failed to serve gRPC server: %v\n", err)
	}
}
