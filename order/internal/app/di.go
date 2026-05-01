package app

import (
	"context"
	"sync"

	orderV1Api "github.com/ArchibaldKronin/microservices_test/order/internal/api/order/v1"
	clientGRPC "github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc"
	inventoryClient "github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc/payment/v1"
	"github.com/ArchibaldKronin/microservices_test/order/internal/config"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	orderRepository "github.com/ArchibaldKronin/microservices_test/order/internal/repository/order"
	"github.com/ArchibaldKronin/microservices_test/order/internal/service"
	orderService "github.com/ArchibaldKronin/microservices_test/order/internal/service/order"
	txmanager "github.com/ArchibaldKronin/microservices_test/order/internal/txManager"
	def "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	orderV1APIHandler def.Handler

	orderService service.OrderService

	connectionInventory *grpc.ClientConn
	inventoryClient     clientGRPC.InventoryClient

	connectionPayment *grpc.ClientConn
	paymentClient     clientGRPC.PaymentClient

	orderTxManager txmanager.TxManager

	orderRepository repository.OrderRepository

	pgPool   *pgxpool.Pool
	pgPoolMu sync.Mutex
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1APIHandler(ctx context.Context) (def.Handler, error) {
	if d.orderV1APIHandler == nil {
		service, err := d.OrderService(ctx)
		if err != nil {
			return nil, err
		}

		d.orderV1APIHandler = orderV1Api.NewApi(service)
	}

	return d.orderV1APIHandler, nil
}

func (d *diContainer) OrderService(ctx context.Context) (service.OrderService, error) {
	if d.orderService == nil {
		var drepo repository.OrderRepository
		var dtxm txmanager.TxManager
		var dinvCl clientGRPC.InventoryClient
		var dpayCl clientGRPC.PaymentClient

		errG, ctx := errgroup.WithContext(ctx)

		errG.Go(func() error {
			repo, err := d.OrderRepository(ctx)
			if err != nil {
				return err
			}
			drepo = repo
			return nil
		})

		errG.Go(func() error {
			txm, err := d.OrderTxManaget(ctx)
			if err != nil {
				return err
			}
			dtxm = txm
			return nil
		})

		errG.Go(func() error {
			invCl, err := d.InventoryClient(ctx)
			if err != nil {
				return err
			}
			dinvCl = invCl
			return nil
		})

		errG.Go(func() error {
			payCl, err := d.PaymentClient(ctx)
			if err != nil {
				return err
			}
			dpayCl = payCl
			return nil
		})

		err := errG.Wait()
		if err != nil {
			return nil, err
		}

		d.orderService = orderService.NewService(drepo, dtxm, dinvCl, dpayCl)
	}

	return d.orderService, nil
}

func (d *diContainer) InventoryClientConnection(_ context.Context) (*grpc.ClientConn, error) {
	if d.connectionInventory == nil {
		connInventory, err := grpc.NewClient(
			config.AppConfig().InventoryGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}
		d.connectionInventory = connInventory
	}

	return d.connectionInventory, nil
}

func (d *diContainer) InventoryClient(ctx context.Context) (clientGRPC.InventoryClient, error) {
	if d.inventoryClient == nil {
		connInventory, err := d.InventoryClientConnection(ctx)
		if err != nil {
			return nil, err
		}

		generatedInventoryCl := inventory_v1.NewInventoryServiceClient(connInventory)
		d.inventoryClient = inventoryClient.NewClient(generatedInventoryCl)
	}

	return d.inventoryClient, nil
}

func (d *diContainer) PaymentClientConnection(_ context.Context) (*grpc.ClientConn, error) {
	if d.connectionPayment == nil {
		connPayment, err := grpc.NewClient(
			config.AppConfig().PaymentGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}
		d.connectionPayment = connPayment
	}

	return d.connectionPayment, nil
}

func (d *diContainer) PaymentClient(ctx context.Context) (clientGRPC.PaymentClient, error) {
	if d.paymentClient == nil {
		connPayment, err := d.PaymentClientConnection(ctx)
		if err != nil {
			return nil, err
		}

		generatedPaymentCl := payment_v1.NewPaymentServiceClient(connPayment)
		d.paymentClient = paymentClient.NewClient(generatedPaymentCl)
	}

	return d.paymentClient, nil
}

func (d *diContainer) OrderTxManaget(ctx context.Context) (txmanager.TxManager, error) {
	if d.orderTxManager == nil {

		pool, err := d.PgPool(ctx)
		if err != nil {
			return nil, err
		}

		d.orderTxManager = txmanager.NewTxRepoManager(pool)
	}

	return d.orderTxManager, nil
}

func (d *diContainer) OrderRepository(ctx context.Context) (repository.OrderRepository, error) {
	if d.orderRepository == nil {

		pool, err := d.PgPool(ctx)
		if err != nil {
			return nil, err
		}

		d.orderRepository = orderRepository.NewRepository(pool)
	}

	return d.orderRepository, nil
}

func (d *diContainer) PgPool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.pgPool != nil {
		return d.pgPool, nil
	}

	d.pgPoolMu.Lock()
	defer d.pgPoolMu.Unlock()

	if d.pgPool != nil {
		return d.pgPool, nil
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			return nil, err
		}

		if err := pool.Ping(ctx); err != nil {
			return nil, err
		}

		d.pgPool = pool
	}

	return d.pgPool, nil
}
