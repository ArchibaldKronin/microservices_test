package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	OTLPAddress() string
	ServiceName() string
	ServiceEnvironment() string
	EnableOTLP() bool
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

type IamGRPCClientConfig interface {
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

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
	ServiceName() string
}
