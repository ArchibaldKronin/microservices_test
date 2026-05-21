package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/config"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/grpc/health"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/tracing"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
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
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initTracer,
		a.initMongo,
		a.initAuthInterceptor,
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

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	err := logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger,
	)
	if err != nil {
		log.Printf("❌ logger init failed: %v; using noop logger", err)
		return err
	}

	closer.AddNamed("Logger", func(ctx context.Context) error {
		_ = logger.Sync()
		_ = logger.Close()
		return nil
	})

	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initMongo(ctx context.Context) error {
	client, err := a.diContainer.MongoDBClient(ctx)
	if err != nil {
		logger.Error(ctx, "❌ Ошибка инициализации клиента Mongo", zap.Error(err))
		return err
	}

	closer.AddNamed("MongoDB client", client.Disconnect)

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

func (a *App) initAuthInterceptor(ctx context.Context) error {
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

	_, err = a.diContainer.AuthInterceptor(ctx)
	if err != nil {
		logger.Error(ctx, "❌ error init auth interceptor", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) initListener(ctx context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
	if err != nil {
		logger.Error(ctx, "failed to listen tcp", zap.String("address", config.AppConfig().InventoryGRPC.Address()), zap.Error(err))
		return err
	}

	closer.AddNamed(
		"TCP listner",
		func(ctx context.Context) error {
			lerr := listener.Close()
			if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
				logger.Error(ctx, "failed to closer listener", zap.String("address", config.AppConfig().InventoryGRPC.Address()), zap.Error(lerr))
				return lerr
			}

			return nil
		},
	)

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	//nolint:gosec
	authInterceptor, _ := a.diContainer.AuthInterceptor(ctx)

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			authInterceptor.Unary(),
			tracing.UnaryServerInterceptor("inventory-service"),
		),
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

	api, err := a.diContainer.InventoryV1API(ctx)
	if err != nil {
		logger.Error(ctx, "failed to create api", zap.Error(err))
	}
	inventory_v1.RegisterInventoryServiceServer(
		a.grpcServer,
		api,
	)

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC InventoryService server listening on %s", config.AppConfig().InventoryGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		logger.Error(ctx, "failed to serve Inventory", zap.Error(err))
		return err
	}

	return nil
}
