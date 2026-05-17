package app

import (
	"context"
	"sync"

	orderV1Api "github.com/ArchibaldKronin/microservices_test/order/internal/api/order/v1"
	clientGRPC "github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc"
	inventoryClient "github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc/payment/v1"
	"github.com/ArchibaldKronin/microservices_test/order/internal/config"
	kafkaConverter "github.com/ArchibaldKronin/microservices_test/order/internal/converter/kafka"
	shipAssembledDecoder "github.com/ArchibaldKronin/microservices_test/order/internal/converter/kafka/decoder"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	orderRepository "github.com/ArchibaldKronin/microservices_test/order/internal/repository/order"
	"github.com/ArchibaldKronin/microservices_test/order/internal/service"
	orderConsumer "github.com/ArchibaldKronin/microservices_test/order/internal/service/consumer/order_consumer"
	orderService "github.com/ArchibaldKronin/microservices_test/order/internal/service/order"
	orderProducer "github.com/ArchibaldKronin/microservices_test/order/internal/service/producer/order_producer"
	txmanager "github.com/ArchibaldKronin/microservices_test/order/internal/txManager"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	wrappedKafka "github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/producer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	authMiddlewarePackage "github.com/ArchibaldKronin/microservices_test/platform/pkg/middleware/http"
	kafkaMiddleware "github.com/ArchibaldKronin/microservices_test/platform/pkg/middleware/kafka"
	def "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
	auth_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/auth/v1"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	orderV1APIHandler def.Handler ////////////////////////////

	orderService         service.OrderService         ////////////////////////////
	orderProducerService service.OrderProducerService ////////////////////////////
	orderConsumerService service.OrderConsumerService ////////////////////////////

	connectionInventory *grpc.ClientConn           ////////////////////////////
	inventoryClient     clientGRPC.InventoryClient ////////////////////////////

	connectionPayment *grpc.ClientConn         ////////////////////////////
	paymentClient     clientGRPC.PaymentClient ////////////////////////////

	consumerGroup          sarama.ConsumerGroup
	orderAssembledConsumer wrappedKafka.Consumer ////////////////////////////

	orderAssembledDecoder kafkaConverter.ShipAssembledDecoder
	syncProducer          sarama.SyncProducer   ////////////////////////////
	orderPaidProducer     wrappedKafka.Producer ////////////////////////////

	orderTxManager txmanager.TxManager ////////////////////////////

	orderRepository repository.OrderRepository ////////////////////////////

	pgPool   *pgxpool.Pool ////////////////////////////
	pgPoolMu sync.Mutex    ////////////////////////////
	/////////////////////////////////////////////////////
	authMiddleware      *authMiddlewarePackage.AuthMiddleware
	iamClient           auth_v1.AuthServiceClient
	iamClientConnection *grpc.ClientConn

	mu sync.Mutex
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthMiddleware(ctx context.Context) (*authMiddlewarePackage.AuthMiddleware, error) {
	if d.authMiddleware == nil {
		client, err := d.IamClient(ctx)
		if err != nil {
			return nil, err
		}

		d.authMiddleware = authMiddlewarePackage.NewAuthMiddleware(client)
	}

	return d.authMiddleware, nil
}

func (d *diContainer) IamClient(ctx context.Context) (auth_v1.AuthServiceClient, error) {
	if d.iamClient == nil {
		conn, err := d.IamClientConnection(ctx)
		if err != nil {
			return nil, err
		}

		d.iamClient = auth_v1.NewAuthServiceClient(conn)
	}

	return d.iamClient, nil
}

func (d *diContainer) IamClientConnection(_ context.Context) (*grpc.ClientConn, error) {
	if d.iamClientConnection == nil {
		connIam, err := grpc.NewClient(
			config.AppConfig().IamConfig.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}
		d.iamClientConnection = connIam
	}

	return d.iamClientConnection, nil
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

func (d *diContainer) OrderConsumerService(ctx context.Context) (service.OrderConsumerService, error) {
	if d.orderConsumerService == nil {

		service, err := d.OrderService(ctx)
		if err != nil {
			return nil, err
		}

		orderAssembledConsumer, err := d.OrderAssembledConsumer(ctx)
		if err != nil {
			return nil, err
		}

		orderAssembledDecoder, err := d.OrderAssembledDecoder(ctx)
		if err != nil {
			return nil, err
		}

		d.orderConsumerService = orderConsumer.NewService(
			orderAssembledConsumer,
			orderAssembledDecoder,
			service,
		)
	}

	return d.orderConsumerService, nil
}

func (d *diContainer) OrderService(ctx context.Context) (service.OrderService, error) {
	if d.orderService != nil {
		return d.orderService, nil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.orderService != nil {
		return d.orderService, nil
	}

	if d.orderService == nil {
		var drepo repository.OrderRepository
		var dtxm txmanager.TxManager
		var dinvCl clientGRPC.InventoryClient
		var dpayCl clientGRPC.PaymentClient
		var dorderProducerService service.OrderProducerService

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

		errG.Go(func() error {
			orderProducerService, err := d.OrderProducerService(ctx)
			if err != nil {
				return err
			}
			dorderProducerService = orderProducerService
			return nil
		})

		err := errG.Wait()
		if err != nil {
			return nil, err
		}

		d.orderService = orderService.NewService(drepo, dtxm, dinvCl, dpayCl, dorderProducerService)
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

func (d *diContainer) OrderProducerService(ctx context.Context) (service.OrderProducerService, error) {
	if d.orderProducerService == nil {
		orderPaidProducer, err := d.OrderPaidProducer(ctx)
		if err != nil {
			return nil, err
		}

		d.orderProducerService = orderProducer.NewService(orderPaidProducer)
	}

	return d.orderProducerService, nil
}

func (d *diContainer) OrderPaidProducer(ctx context.Context) (wrappedKafka.Producer, error) {
	if d.orderPaidProducer == nil {

		syncProd, err := d.SyncProducer(ctx)
		if err != nil {
			return nil, err
		}

		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			syncProd,
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderPaidProducer, nil
}

func (d *diContainer) SyncProducer(ctx context.Context) (sarama.SyncProducer, error) {
	if d.syncProducer == nil {
		syncProd, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
		)
		if err != nil {
			return nil, err
		}

		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return syncProd.Close()
		})

		d.syncProducer = syncProd
	}

	return d.syncProducer, nil
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

func (d *diContainer) OrderAssembledConsumer(ctx context.Context) (wrappedKafka.Consumer, error) {
	if d.orderAssembledConsumer == nil {

		consumerGroup, err := d.ConsumerGroup(ctx)
		if err != nil {
			return nil, err
		}
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			consumerGroup,
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer, nil
}

func (d *diContainer) ConsumerGroup(ctx context.Context) (sarama.ConsumerGroup, error) {
	if d.consumerGroup == nil {

		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			return nil, err
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup, nil
}

func (d *diContainer) OrderAssembledDecoder(ctx context.Context) (kafkaConverter.ShipAssembledDecoder, error) {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = shipAssembledDecoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder, nil
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
