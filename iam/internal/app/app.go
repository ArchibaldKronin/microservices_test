package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/config"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/grpc/health"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	pgMigrator "github.com/ArchibaldKronin/microservices_test/platform/pkg/migrator/pg"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/tracing"
	auth_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/auth/v1"
	user_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/user/v1"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listner     net.Listener
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
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initLogger,
		a.initDI,
		a.initCloser,
		a.initTracer,
		a.initPostgres,
		a.initMigrator,
		a.initRedis,
		a.initListener,
		a.initGRPCServer,
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

func (a *App) initDI(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initCloser(ctx context.Context) error {
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

func (a *App) initPostgres(ctx context.Context) error {
	pool, err := a.diContainer.UserPgPool(ctx)
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
	pool, err := a.diContainer.UserPgPool(ctx)
	if err != nil {
		logger.Error(ctx, "❌ Ошибка инициализации pgxpool", zap.Error(err))
		return err
	}

	migratorRunner := pgMigrator.NewMigrator(
		stdlib.OpenDBFromPool(pool),
		config.AppConfig().Postgres.MigrationsDir(),
	)
	ctxMigrator, cancelMigrator := context.WithTimeout(ctx, 30*time.Second)
	defer cancelMigrator()

	err = migratorRunner.Up(ctxMigrator)
	if err != nil {
		logger.Error(ctxMigrator, "❌ databse migration error", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) initRedis(ctx context.Context) error {
	pool, err := a.diContainer.AuthRedisPool(ctx)
	if err != nil {
		logger.Error(ctx, "❌ Ошибка инициализации Redis pool", zap.Error(err))
		return err
	}

	closer.AddNamed("Redis Pool Connection", func(ctx context.Context) error {
		err = pool.Close()
		if err != nil {
			logger.Error(ctx, "❌ Ошибка закрытия Redis pool", zap.Error(err))
			return err
		}
		return nil
	})

	return nil
}

func (a *App) initListener(ctx context.Context) error {
	listner, err := net.Listen("tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		logger.Error(ctx, "failed to listen tcp", zap.String("address", config.AppConfig().GRPC.Address()), zap.Error(err))
		return err
	}

	closer.AddNamed(
		"TCP Listner",
		func(ctx context.Context) error {
			lerr := listner.Close()
			if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
				logger.Error(ctx, "failed to closer listener", zap.String("address", config.AppConfig().GRPC.Address()), zap.Error(lerr))
				return lerr
			}

			return nil
		},
	)

	a.listner = listner

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(tracing.UnaryServerInterceptor("iam-service")),
	)

	closer.AddNamed(
		"gRPC server",
		func(ctx context.Context) error {
			a.grpcServer.GracefulStop()
			return nil
		},
	)

	reflection.Register(a.grpcServer)

	health.RegisterService(a.grpcServer)

	userApi, err := a.diContainer.UserV1Api(ctx)
	if err != nil {
		logger.Error(ctx, "failed to create User api", zap.Error(err))
	}
	authApi, err := a.diContainer.AuthV1Api(ctx)
	if err != nil {
		logger.Error(ctx, "failed to create Auth api", zap.Error(err))
	}

	user_v1.RegisterUserServiceServer(a.grpcServer, userApi)
	auth_v1.RegisterAuthServiceServer(a.grpcServer, authApi)

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC IAMService server listening on %s", config.AppConfig().GRPC.Address()))

	err := a.grpcServer.Serve(a.listner)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		logger.Error(ctx, "failed to serve Inventory", zap.Error(err))
		return err
	}

	return nil
}
