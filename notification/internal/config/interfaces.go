package config

import "github.com/IBM/sarama"

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

type ConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type TelegramConfig interface {
	Token() string
}

// type OrderPaidConsumerConfig interface {
// 	Topic() string
// 	GroupID() string
// 	Config() *sarama.Config
// }

// type OrderAssembledConsumerConfig interface {
// 	Topic() string
// 	GroupID() string
// 	Config() *sarama.Config
// }
