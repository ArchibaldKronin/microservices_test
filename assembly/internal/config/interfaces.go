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

type KafkaConfig interface {
	Brokers() []string
}

type ShipAssembledProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
	ServiceName() string
}
