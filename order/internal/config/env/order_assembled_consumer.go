package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderAssembledConsumerEnvConfig struct {
	Topic   string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type orderAssembledConsumerConfig struct {
	config orderAssembledConsumerEnvConfig
}

func NewOrderAssembledConsumerConfig() (*orderAssembledConsumerConfig, error) {
	var config orderAssembledConsumerEnvConfig
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &orderAssembledConsumerConfig{config: config}, nil
}

func (cfg *orderAssembledConsumerConfig) Topic() string {
	return cfg.config.Topic
}

func (cfg *orderAssembledConsumerConfig) GroupID() string {
	return cfg.config.GroupID
}

func (cfg *orderAssembledConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
