package main

import (
	"log"
	"net"
	"os"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/jany/my-coffee/config"
	brewpb "github.com/jany/my-coffee/gen/proto/brew"
	"github.com/jany/my-coffee/internal/brews"
	database "github.com/jany/my-coffee/internal/datbase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load config
	config.Load()

	// Connect to database
	db := database.Connect()
	defer database.Close()

	// Create gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed to create validator: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(protovalidate_middleware.UnaryServerInterceptor(validator)),
	)
	brewpb.RegisterBrewServiceServer(grpcServer, brews.New(db))

	// Only enable reflection in development
    if os.Getenv("ENV") != "production" {
      reflection.Register(grpcServer)
    }

	log.Println("BrewService running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
