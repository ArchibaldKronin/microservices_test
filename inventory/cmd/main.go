package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	inventoryV1API "github.com/ArchibaldKronin/microservices_test/inventory/internal/api/inventory/v1"
	inventoryRepo "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/part"
	inventoryService "github.com/ArchibaldKronin/microservices_test/inventory/internal/service/part"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50051
	dbURI    = "mongodb://inventory-service-user:inventory-service-password@localhost:27017/inventory-service?authSource=admin"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listenner: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	// repo := inventoryRepo.NewRepository(inventoryRepo.InitialParts)
	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("failed to conncet to DB: %v\n", err)
		return
	}
	defer func() {
		cerr := mongoClient.Disconnect(ctx)
		if cerr != nil {
			log.Printf("failed to disconnect BD: %v\n", cerr)
		}
	}()

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping DB: %v\n", err)
		return
	}

	db := mongoClient.Database("inventory-service")

	repo, err := inventoryRepo.NewRepository(db)
	service := inventoryService.NewService(repo)
	api := inventoryV1API.NewApi(service)

	inventory_v1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC Inventory server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve Inventory: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC Inventory server...")
	s.GracefulStop()
	log.Println("✅ Inventory server stopped")
}
