package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/app"
	"github.com/ArchibaldKronin/microservices_test/iam/internal/config"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

const (
	configPath = "../deploy/compose/iam/.env"
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

	initCtx, initCancel := context.WithTimeout(appCtx, 7*time.Second)
	defer initCancel()

	app, err := app.New(initCtx)
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
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "error shutting down", zap.Error(err))
	}
}
