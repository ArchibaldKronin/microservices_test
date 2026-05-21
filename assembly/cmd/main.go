package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArchibaldKronin/microservices_test/assembly/internal/app"
	"github.com/ArchibaldKronin/microservices_test/assembly/internal/config"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

const (
	configPath = "../deploy/compose/assembly/.env"
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

	// closer.Configure(
	// 	syscall.SIGINT,
	// 	syscall.SIGTERM,
	// )

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

	<-appCtx.Done()
	logger.Info(appCtx, "🛑 Получен системный сигнал, начинаем graceful shutdown...")
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "error shitting down", zap.Error(err))
	}
}
