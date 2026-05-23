package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ArchibaldKronin/microservices_test/order/internal/api/health"
	"github.com/ArchibaldKronin/microservices_test/order/internal/config"
	"github.com/ArchibaldKronin/microservices_test/order/internal/metrics"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	metricsProvider "github.com/ArchibaldKronin/microservices_test/platform/pkg/metrics"
	pgMigrator "github.com/ArchibaldKronin/microservices_test/platform/pkg/migrator/pg"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/tracing"
	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// закрывать клиент инвентори
// закрывать клиент пэймент
// закрывать pgpool
// init migrator
const (
	requestProcessingTimeout = 10 * time.Second
	shutdownTimeout          = 5 * time.Second
)

type App struct {
	diContainer *diContainer
	orderServer *order_v1.Server
	httpServer  *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer crashed: %w", err)
		}
	}()

	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- fmt.Errorf("HTTP server crashed: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		logger.Error(ctx, "Component crushed, shutting down", zap.Error(err))
		cancel()
		<-ctx.Done()
		return err
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	}

	return nil
	// return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initLogger,
		a.initDi,
		a.initCloser,
		a.initTracer,
		a.initMetricsProvider,
		a.initMetrics,
		a.initPg,
		a.initMigrator,
		a.initInventoryConnection,
		a.initPaymentConnection,
		a.initAuthMiddleware,
		a.initOrderServer,
		a.initHTTPServer,
	}

	for _, f := range inits {
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

func (a *App) initDi(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initTracer(ctx context.Context) error {
	err := tracing.InitTracer(ctx, config.AppConfig().Tracing)
	if err != nil {
		logger.Error(ctx, "❌ Ошибка инициализации OTEL tracer", zap.Error(err))
		return err
	}

	closer.AddNamed("tracer", tracing.ShutdownTracer)

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

func (a *App) initPg(ctx context.Context) error {
	pool, err := a.diContainer.PgPool(ctx)
	if err != nil {
		logger.Error(ctx, "❌ Ошибка инициализации pgxpool", zap.Error(err))
		return err
	}

	closer.AddNamed("PgPoolConnection", func(ctx context.Context) error {
		pool.Close()
		return nil
	})

	return nil
}

func (a *App) initMigrator(ctx context.Context) error {
	pool, err := a.diContainer.PgPool(ctx)
	if err != nil {
		logger.Error(ctx, "❌ Ошибка инициализации pgxpool", zap.Error(err))
		return err
	}

	migratorRunner := pgMigrator.NewMigrator(stdlib.OpenDBFromPool(pool), config.AppConfig().Postgres.MigrationsDir())
	ctxMigrator, cancelMigrator := context.WithTimeout(ctx, 30*time.Second)
	defer cancelMigrator()

	err = migratorRunner.Up(ctxMigrator)
	if err != nil {
		logger.Error(ctxMigrator, "❌ databse migration error", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) initInventoryConnection(ctx context.Context) error {
	clientConn, err := a.diContainer.InventoryClientConnection(ctx)
	if err != nil {
		logger.Error(ctx, "❌ error init inventory client", zap.Error(err))
		return err
	}

	closer.AddNamed(
		"Inventory connection",
		func(ctx context.Context) error {
			if cerr := clientConn.Close(); cerr != nil {
				logger.Error(ctx, "❌ failed to close connect to Inventory", zap.Error(cerr))
				return cerr
			}

			return nil
		},
	)

	return nil
}

func (a *App) initPaymentConnection(ctx context.Context) error {
	clientConn, err := a.diContainer.PaymentClientConnection(ctx)
	if err != nil {
		logger.Error(ctx, "❌ error init payment client", zap.Error(err))
		return err
	}

	closer.AddNamed(
		"Payment connection",
		func(ctx context.Context) error {
			if cerr := clientConn.Close(); cerr != nil {
				logger.Error(ctx, "❌ failed to close connect to Payment", zap.Error(cerr))
				return cerr
			}

			return nil
		},
	)

	return nil
}

func (a *App) initAuthMiddleware(ctx context.Context) error {
	conn, err := a.diContainer.IamClientConnection(ctx)
	if err != nil {
		logger.Error(ctx, "❌ error init iam client", zap.Error(err))
		return err
	}

	closer.AddNamed(
		"Iam connection",
		func(ctx context.Context) error {
			if cerr := conn.Close(); cerr != nil {
				logger.Error(ctx, "❌ failed to close connect to Iam", zap.Error(cerr))
				return cerr
			}

			return nil
		},
	)

	_, err = a.diContainer.AuthMiddleware(ctx)
	if err != nil {
		logger.Error(ctx, "❌ error init auth middleware", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) initOrderServer(ctx context.Context) error {
	handler, err := a.diContainer.OrderV1APIHandler(ctx)
	if err != nil {
		logger.Error(ctx, "❌ ошибка создания Order сервера OpenAPI", zap.Error(err))
		return err
	}

	orderServer, err := order_v1.NewServer(handler)
	if err != nil {
		logger.Error(ctx, "❌ ошибка создания Order сервера OpenAPI", zap.Error(err))
		return err
	}

	a.orderServer = orderServer
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(requestProcessingTimeout))
	r.Use(tracing.HTTPHandlerMiddleware("order-service"))
	// r.Use(a.diContainer.authMiddleware.Handle)

	r.Get("/health", health.HealthCheck)
	r.Get("/ready", health.ReadyCheck(
		a.diContainer.pgPool,
		a.diContainer.connectionInventory,
		a.diContainer.connectionPayment,
	))

	r.Group(func(protected chi.Router) {
		protected.Use(a.diContainer.authMiddleware.Handle)
		protected.Mount("/", a.orderServer)
	})

	// r.Mount("/", a.orderServer)

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTime(),
	}

	closer.AddNamed(
		"HTTP server",
		func(ctx context.Context) error {
			ctxShutDown, cancel := context.WithTimeout(ctx, shutdownTimeout)
			defer cancel()

			return a.httpServer.Shutdown(ctxShutDown)
		},
	)

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP-сервер запущен на порту: %s", config.AppConfig().OrderHTTP.Port()))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, "❌ Ошибка работы сервера", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderPaid Kafka consumer running")

	c, err := a.diContainer.OrderConsumerService(ctx)
	if err != nil {
		return err
	}

	err = c.RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
