package app

import (
	"context"
	"sync"
	"time"

	clients "github.com/ArchibaldKronin/microservices_test/notification/internal/client/http"
	telegramClient "github.com/ArchibaldKronin/microservices_test/notification/internal/client/http/telegram"
	"github.com/ArchibaldKronin/microservices_test/notification/internal/config"
	kafkaConverter "github.com/ArchibaldKronin/microservices_test/notification/internal/converter/kafka"
	order_Decoder "github.com/ArchibaldKronin/microservices_test/notification/internal/converter/kafka/decoder"
	"github.com/ArchibaldKronin/microservices_test/notification/internal/service"
	"github.com/ArchibaldKronin/microservices_test/notification/internal/service/consumer/order_assembled_consumer"
	"github.com/ArchibaldKronin/microservices_test/notification/internal/service/consumer/order_paid_consumer"
	"github.com/ArchibaldKronin/microservices_test/notification/internal/service/telegram"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/closer"
	wrappedKafka "github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/consumer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	kafkaMiddleware "github.com/ArchibaldKronin/microservices_test/platform/pkg/middleware/kafka"
	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

type diContainer struct {
	orderPaidConsumerService      service.OrderPaidConsumerService
	orderAssembledConsumerService service.OrderAssembledConsumerService
	telegramService               service.TelegramService

	telegramClient clients.TelegramClient
	telegramBot    *bot.Bot

	orderPaidConsumer      wrappedKafka.Consumer
	consumerGroupOrderPaid sarama.ConsumerGroup
	orderPaidDecoder       kafkaConverter.OrderPaidDecoder

	orderAssembledConsumer      wrappedKafka.Consumer
	consumerGroupOrderAssembled sarama.ConsumerGroup
	orderAssembledDecoder       kafkaConverter.OrderAssembledDecoder

	telegramServiceOnce sync.Once
	telegramServiceErr  error
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) (service.OrderAssembledConsumerService, error) {
	if d.orderAssembledConsumerService == nil {
		orderAssembledConsumer, err := d.OrderAssembledConsumer(ctx)
		if err != nil {
			return nil, err
		}

		orderAssembledDecoder, err := d.OrderAssembledDecoder(ctx)
		if err != nil {
			return nil, err
		}

		telegramService, err := d.TelegramService(ctx)
		if err != nil {
			return nil, err
		}

		d.orderAssembledConsumerService = order_assembled_consumer.NewService(
			orderAssembledConsumer,
			orderAssembledDecoder,
			telegramService,
		)
	}

	return d.orderAssembledConsumerService, nil
}

func (d *diContainer) OrderAssembledConsumer(ctx context.Context) (wrappedKafka.Consumer, error) {
	if d.orderAssembledConsumer == nil {
		consumerGroup, err := d.ConsumerGroupOrderAssembled(ctx)
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

func (d *diContainer) ConsumerGroupOrderAssembled(ctx context.Context) (sarama.ConsumerGroup, error) {
	if d.consumerGroupOrderAssembled == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			return nil, err
		}

		closer.AddNamed("Kafka consumer group OrderAssembled", func(ctx context.Context) error {
			return d.consumerGroupOrderAssembled.Close()
		})

		d.consumerGroupOrderAssembled = consumerGroup
	}

	return d.consumerGroupOrderAssembled, nil
}

func (d *diContainer) OrderAssembledDecoder(ctx context.Context) (kafkaConverter.OrderAssembledDecoder, error) {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = order_Decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder, nil
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) (service.OrderPaidConsumerService, error) {
	if d.orderPaidConsumerService == nil {
		orderPaidConsumer, err := d.OrderPaidConsumer(ctx)
		if err != nil {
			return nil, err
		}

		orderPaidDecoder, err := d.OrderPaidDecoder(ctx)
		if err != nil {
			return nil, err
		}

		telegramService, err := d.TelegramService(ctx)
		if err != nil {
			return nil, err
		}

		d.orderPaidConsumerService = order_paid_consumer.NewService(
			orderPaidConsumer,
			orderPaidDecoder,
			telegramService,
		)
	}

	return d.orderPaidConsumerService, nil
}

func (d *diContainer) OrderPaidConsumer(ctx context.Context) (wrappedKafka.Consumer, error) {
	if d.orderPaidConsumer == nil {
		consumerGroup, err := d.ConsumerGroupOrderPaid(ctx)
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

func (d *diContainer) ConsumerGroupOrderPaid(ctx context.Context) (sarama.ConsumerGroup, error) {
	if d.consumerGroupOrderPaid == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			return nil, err
		}

		closer.AddNamed("Kafka consumer group OrderPaid", func(ctx context.Context) error {
			return d.consumerGroupOrderPaid.Close()
		})

		d.consumerGroupOrderPaid = consumerGroup
	}

	return d.consumerGroupOrderPaid, nil
}

func (d *diContainer) OrderPaidDecoder(ctx context.Context) (kafkaConverter.OrderPaidDecoder, error) {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = order_Decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder, nil
}

func (d *diContainer) TelegramService(ctx context.Context) (service.TelegramService, error) {
	d.telegramServiceOnce.Do(func() {
		tgClient, innerErr := d.TelegramClient(ctx)
		if innerErr != nil {
			d.telegramServiceErr = innerErr
			return
		}

		d.telegramService = telegram.NewService(tgClient)
	})

	return d.telegramService, d.telegramServiceErr
}

func (d *diContainer) TelegramClient(ctx context.Context) (clients.TelegramClient, error) {
	if d.telegramClient == nil {
		tgBot, err := d.TelegramBot(ctx)
		if err != nil {
			return nil, err
		}

		d.telegramClient = telegramClient.NewTelegramClient(tgBot)
	}

	return d.telegramClient, nil
}

func (d *diContainer) TelegramBot(ctx context.Context) (*bot.Bot, error) {
	if d.telegramBot == nil {

		const (
			maxAttempts = 3
			retryDelay  = 3 * time.Second
		)

		var tgBot *bot.Bot
		var err error
		for i := 0; i < maxAttempts; i++ {
			tgBot, err = bot.New(config.AppConfig().Telegram.Token())
			if err == nil {
				break
			}

			if i == maxAttempts-1 {
				logger.Error(ctx, "Failed to connect Telegram API after 3 times", zap.Error(err))
				return nil, err
			}

			logger.Warn(ctx, "Failed to connect Telegram API", zap.Int("attempt", i+1), zap.Error(err))
			select {
			case <-time.After(retryDelay):
				logger.Debug(ctx, "SELECT TIME AFTER")
				continue
			case <-ctx.Done():
				logger.Debug(ctx, "SELECT CTX DONE")
				return nil, ctx.Err()
			}
		}

		d.telegramBot = tgBot
	}

	return d.telegramBot, nil
}
