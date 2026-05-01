package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArchibaldKronin/microservices_test/order/internal/app"
	"github.com/ArchibaldKronin/microservices_test/order/internal/config"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

const (
	// httpPort                 = "8080"
	// dbURI                    = "postgres://order-service-user:order-service-password@localhost:5432/order-service"
	// readHeaderTimeout        = 5 * time.Second
	// requestProcessingTimeout = 10 * time.Second
	// shutdownTimeout          = 10 * time.Second
	// migrationsDir            = "./migrations"
	configPath = "../deploy/compose/order/.env"
)

const (
// serverInventoryAddress = "localhost:50051"
// serverPaymentAddress   = "localhost:50052"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	app, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = app.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "error shutting down", zap.Error(err))
	}
}
