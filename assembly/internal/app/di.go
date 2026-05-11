package app

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/assembly/internal/config"
	kafkaConverter "github.com/ArchibaldKronin/microservices_test/assembly/internal/converter/kafka"
	orderPaidDecoder "github.com/ArchibaldKronin/microservices_test/assembly/internal/converter/kafka/decoder"
	"github.com/ArchibaldKronin/microservices_test/assembly/internal/service"
	orderConsumer "github.com/ArchibaldKronin/microservices_test/assembly/internal/service/consumer/order_consumer"
	orderProducer "github.com/ArchibaldKronin/microservices_test/assembly/internal/service/producer/order_producer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	wrappedKafka "github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/producer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	kafkaMiddleware "github.com/ArchibaldKronin/microservices_test/platform/pkg/middleware/kafka"
	"github.com/IBM/sarama"
)

type diContainer struct {
	shipAssembledProducerService service.ShipAssembledProducer ////////////////////
	orderPaidConsumerService     service.OrderPaidConsumer

	orderPaidConsumer wrappedKafka.Consumer
	consumerGroup     sarama.ConsumerGroup
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder

	syncProducer          sarama.SyncProducer   ////////////////////
	shipAssembledProducer wrappedKafka.Producer ////////////////////
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) (service.OrderPaidConsumer, error) {
	if d.orderPaidConsumerService == nil {
		orderPaidConsumer, err := d.OrderPaidConsumer(ctx)
		if err != nil {
			return nil, err
		}
		orderPaidDecoder, err := d.OrderPaidDecoder(ctx)
		if err != nil {
			return nil, err
		}
		shipAssembledProducerService, err := d.ShipAssembledProducerService(ctx)
		if err != nil {
			return nil, err
		}

		d.orderPaidConsumerService = orderConsumer.NewService(
			orderPaidConsumer,
			orderPaidDecoder,
			shipAssembledProducerService,
		)
	}

	return d.orderPaidConsumerService, nil
}

func (d *diContainer) OrderPaidConsumer(ctx context.Context) (wrappedKafka.Consumer, error) {
	if d.orderPaidConsumer == nil {
		consumerGroup, err := d.ConsumerGroup(ctx)
		if err != nil {
			return nil, err
		}

		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			consumerGroup,
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer, nil
}

func (d *diContainer) ConsumerGroup(ctx context.Context) (sarama.ConsumerGroup, error) {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
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

func (d *diContainer) OrderPaidDecoder(ctx context.Context) (kafkaConverter.OrderPaidDecoder, error) {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = orderPaidDecoder.NewOrderpaidDecoder()
	}

	return d.orderPaidDecoder, nil
}

func (d *diContainer) ShipAssembledProducerService(ctx context.Context) (service.ShipAssembledProducer, error) {
	if d.shipAssembledProducerService == nil {
		shipAssembledProducer, err := d.ShipAssembledProducer(ctx)
		if err != nil {
			return nil, err
		}

		d.shipAssembledProducerService = orderProducer.NewService(shipAssembledProducer)
	}

	return d.shipAssembledProducerService, nil
}

func (d *diContainer) ShipAssembledProducer(ctx context.Context) (wrappedKafka.Producer, error) {
	if d.shipAssembledProducer == nil {
		syncProd, err := d.SyncProducer(ctx)
		if err != nil {
			return nil, err
		}

		d.shipAssembledProducer = wrappedKafkaProducer.NewProducer(
			syncProd,
			config.AppConfig().ShipAssembledProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.shipAssembledProducer, nil
}

func (d *diContainer) SyncProducer(ctx context.Context) (sarama.SyncProducer, error) {
	if d.syncProducer == nil {
		syncProd, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().ShipAssembledProducer.Config(),
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
