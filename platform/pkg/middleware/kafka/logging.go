package kafka

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/consumer"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
}

func Logging(logger Logger) consumer.Middleware {
	return func(next consumer.MessageHandler) consumer.MessageHandler {
		return func(ctx context.Context, msg consumer.Message) error {
			logger.Info(ctx, "Kafka message received", zap.String("topic", msg.Topic))
			return next(ctx, msg)
		}
	}
}
