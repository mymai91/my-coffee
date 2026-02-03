package main

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	menupb "github.com/jany/my-coffee/gen/proto/menu"
	"github.com/jany/my-coffee/internal/menus"
)

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		fmt.Printf("Failed to listen on Port 50052: %v\n", err)
		return
	}

	server := grpc.NewServer()
	menuService := menus.New()
	menupb.RegisterMenuServiceServer(server, menuService)

	// Only enable reflection in development
    if os.Getenv("ENV") != "production" {
        reflection.Register(server)
    }

	fmt.Println("Menu Service is running on port 50052...")
	if err := server.Serve(listener); err != nil {
		fmt.Printf("Failed to serve gRPC server: %v\n", err)
	}
}