package kafka

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, handler consumer.MessageHandler) error
}

type Producer interface {
	Send(ctx context.Context, key, valuy []byte) error
}
