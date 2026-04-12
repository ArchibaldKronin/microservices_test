package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/ArchibaldKronin/microservices_test/order/internal/api/order/v1"
	inventoryClient "github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc/payment/v1"
	"github.com/ArchibaldKronin/microservices_test/order/internal/migrator"
	repo "github.com/ArchibaldKronin/microservices_test/order/internal/repository/order"
	service "github.com/ArchibaldKronin/microservices_test/order/internal/service/order"
	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort                 = "8080"
	dbURI                    = "postgres://order-service-user:order-service-password@localhost:5432/order-service"
	readHeaderTimeout        = 5 * time.Second
	requestProcessingTimeout = 5 * time.Second
	shutdownTimeout          = 10 * time.Second
	migrationsDir            = "../migrations"
)

const (
	serverInventoryAddress = "localhost:50051"
	serverPaymentAddress   = "localhost:50052"
)

func main() {
	/////////////////////////////////////////////////////////////////////// CLIENTS
	connInventory, err := grpc.NewClient(
		serverInventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("ошибка соединения с Inventory: %v", err)
	}
	defer func() {
		if cerr := connInventory.Close(); cerr != nil {
			log.Printf("failed to close connect to Inventory: %v", cerr)
		}
	}()

	generatedInventoryCl := inventory_v1.NewInventoryServiceClient(connInventory)
	inventoryClient := inventoryClient.NewClient(generatedInventoryCl)

	connPayment, err := grpc.NewClient(
		serverPaymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("ошибка соединения с Payment: %v", err)
	}
	defer func() {
		if cerr := connPayment.Close(); cerr != nil {
			log.Printf("failed to close connect to Payment: %v", cerr)
		}
	}()

	generatedPaymentCl := payment_v1.NewPaymentServiceClient(connPayment)
	paymentClient := paymentClient.NewClient(generatedPaymentCl)
	/////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////// SERVER
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
	}
	defer pool.Close()

	migratorRunner := migrator.NewMigrator(stdlib.OpenDBFromPool(pool), migrationsDir)
	err = migratorRunner.Up()
	if err != nil {
		slog.Error("databse migration error", "error ", err)
		return
	}

	storage := repo.NewRepository(pool)

	service := service.NewService(storage, inventoryClient, paymentClient)

	handler := api.NewApi(service)

	orderServer, err := order_v1.NewServer(handler)
	if err != nil {
		log.Printf("ошибка создания Order сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(requestProcessingTimeout))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		///////////////////////////
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		///////////////////////////
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctxShutDown)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
