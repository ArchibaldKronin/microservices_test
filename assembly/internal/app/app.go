package app

import (
	"context"
	"log"

	"github.com/ArchibaldKronin/microservices_test/assembly/internal/config"
	"github.com/ArchibaldKronin/microservices_test/assembly/internal/service/metrics"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	metricsProvider "github.com/ArchibaldKronin/microservices_test/platform/pkg/metrics"
	"go.uber.org/zap"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runConsumer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initLogger,
		a.initDi,
		a.initCloser,
		a.initMetricsProvider,
		a.initMetrics,
	}

	for _, f := range inits {
		f := f
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	err := logger.Init(ctx,
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger,
	)
	if err != nil {
		log.Printf("❌ logger init failed: %v; using noop logger", err)
		logger.SetNopLogger()
		return err
	}

	closer.AddNamed("Logger", func(ctx context.Context) error {
		_ = logger.Sync()     //nolint:gosec
		_ = logger.Close(ctx) //nolint:gosec
		return nil
	})

	return nil
}

func (a *App) initDi(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initMetricsProvider(ctx context.Context) error {
	err := metricsProvider.InitProvider(ctx, config.AppConfig().Metrics)
	if err != nil {
		logger.Error(ctx, "❌ Ошибка инициализации OTEL metrics provider", zap.Error(err))
		return err
	}

	closer.AddNamed("metrics", metricsProvider.Shutdown)

	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	err := metrics.InitMetrics()
	if err != nil {
		logger.Error(ctx, "❌ Ошибка инициализации metrics", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderPaid Kafka consumer running")

	c, err := a.diContainer.OrderPaidConsumerService(ctx)
	if err != nil {
		return err
	}

	err = c.RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
