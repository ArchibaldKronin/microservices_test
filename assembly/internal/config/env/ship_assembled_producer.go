package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type shipAssembledProducerEnvConfig struct {
	Topic string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
}

type shipAssembledProducerConfig struct {
	config shipAssembledProducerEnvConfig
}

func NewshipAssembledProducerConfig() (*shipAssembledProducerConfig, error) {
	var config shipAssembledProducerEnvConfig
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &shipAssembledProducerConfig{config: config}, nil
}

func (cfg *shipAssembledProducerConfig) Topic() string {
	return cfg.config.Topic
}

func (cfg *shipAssembledProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
