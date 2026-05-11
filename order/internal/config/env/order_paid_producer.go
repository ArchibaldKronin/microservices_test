package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderPaidProducerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
}

type orderPaidProducerConfig struct {
	config orderPaidProducerEnvConfig
}

func NewOrderPaidProducerConfig() (*orderPaidProducerConfig, error) {
	var config orderPaidProducerEnvConfig
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &orderPaidProducerConfig{config: config}, nil
}

func (cfg *orderPaidProducerConfig) Topic() string {
	return cfg.config.TopicName
}

func (cfg *orderPaidProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
