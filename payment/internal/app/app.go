package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/ArchibaldKronin/microservices_test/payment/internal/config"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/grpc/health"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
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

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initListner,
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

func (a *App) initDi(_ context.Context) error {
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

func (a *App) initListner(ctx context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
	if err != nil {
		logger.Error(ctx, "failed to listen tcp", zap.String("address", config.AppConfig().PaymentGRPC.Address()), zap.Error(err))
		return err
	}

	closer.AddNamed(
		"TCP listner",
		func(ctx context.Context) error {
			lerr := listener.Close()
			if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
				logger.Error(ctx, "failed to closer listener", zap.String("address", config.AppConfig().PaymentGRPC.Address()), zap.Error(lerr))
				return lerr
			}
			return nil
		},
	)

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(contextEnricher()),
			grpc.UnaryServerInterceptor(loggerUUID()),
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

	api := a.diContainer.PaymentV1API(ctx)

	payment_v1.RegisterPaymentServiceServer(
		a.grpcServer,
		api,
	)

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC Payment service server listening on %s", config.AppConfig().PaymentGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		logger.Error(ctx, "failed to serve Payment", zap.String("address", config.AppConfig().PaymentGRPC.Address()), zap.Error(err))
		return err
	}

	return nil
}

func contextEnricher() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if info.FullMethod == "/payment.v1.PaymentService/PayOrder" {
			if v, ok := req.(*payment_v1.PayOrderRequest); ok {
				ctx = logger.WithUserID(ctx, v.UserUuid)
			}
		}

		return handler(ctx, req)
	}
}

func loggerUUID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, err
		}
		if v, ok := resp.(*payment_v1.PayOrderResponse); ok {
			logger.Info(ctx, "payment postprocessing log", zap.String("transaction_uuid", v.TransactionUuid))
		}
		return resp, err
	}
}
