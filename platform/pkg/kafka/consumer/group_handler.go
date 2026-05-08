package consumer

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type MessageHandler func(ctx context.Context, msg Message) error

type Middleware func(next MessageHandler) MessageHandler

type groupHandler struct {
	handler MessageHandler
	logger  Logger
}

func NewGroupHandler(handler MessageHandler, logger Logger, middlewares ...Middleware) *groupHandler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return &groupHandler{
		handler: handler,
		logger:  logger,
	}
}

func (g *groupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				g.logger.Info(session.Context(), "Kafka message channel closed")
				return nil
			}

			msg := Message{
				Key:            message.Key,
				Value:          message.Value,
				Topic:          message.Topic,
				Partition:      message.Partition,
				Offset:         message.Offset,
				Timestamp:      message.Timestamp,
				BlockTimestamp: message.BlockTimestamp,
				Headers:        extractHeaders(message.Headers),
			}

			var err error

			for i := 0; i < 5; i++ {
				ctx, cancel := context.WithTimeout(session.Context(), 5*time.Second)
				err = g.handler(ctx, msg)
				cancel()

				if err == nil {
					break
				}

				if i == 4 {
					g.logger.Error(session.Context(), "Kafka handler error after 5 retries", zap.Error(err))
					break
				}

				base := time.Second * time.Duration(1<<i)
				// baseBigInt, err := rand.Int(rand.Reader, big.NewInt(1<<i))
				// if err != nil {
				// 	return err
				// }
				// base := time.Second * time.Duration(baseBigInt.Int64())
				// jitterBigInt, err := rand.Int(rand.Reader, big.NewInt(int64(base/2)))
				jitter := time.Duration(rand.Int64N(int64(base / 2))) //nolint:gosec
				// jitter := time.Duration(jitterBigInt.Int64())
				if err != nil {
					return err
				}
				delay := base + jitter
				g.logger.Error(
					session.Context(),
					fmt.Sprintf("Kafka handler error, retry #%d after %s", i+1, delay.String()),
					zap.Error(err),
				)

				select {
				case <-time.After(delay):
					continue
				case <-session.Context().Done():
					g.logger.Info(session.Context(), "Kafka session context done")
					return nil
				}
			}

			if err != nil {
				session.MarkMessage(message, "message wasn't handled")
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			g.logger.Info(session.Context(), "Kafka session context done")
			return nil
		}
	}
}

func extractHeaders(headers []*sarama.RecordHeader) map[string][]byte {
	h := make(map[string][]byte)
	for _, v := range headers {
		if v != nil && v.Key != nil {
			h[string(v.Key)] = v.Value
		}
	}

	return h
}
