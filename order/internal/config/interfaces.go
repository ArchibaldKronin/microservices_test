package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type OrderHTTPConfigConfig interface {
	Address() string
	ReadTime() time.Duration
	Port() string
}

type InventoryGRPCClientConfig interface {
	Address() string
}

type PaymentGRPCClientConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
