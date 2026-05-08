package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderPaidConsumerEnvConfig struct {
	Topic   string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type orderPaidConsumerConfig struct {
	config orderPaidConsumerEnvConfig
}

func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	var config orderPaidConsumerEnvConfig
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &orderPaidConsumerConfig{config: config}, nil
}

func (cfg *orderPaidConsumerConfig) Topic() string {
	return cfg.config.Topic
}

func (cfg *orderPaidConsumerConfig) GroupID() string {
	return cfg.config.GroupID
}

func (cfg *orderPaidConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}

// type OrderPaidConsumerConfig interface {
// 	Topic() string
// 	GroupID() string
// 	Config() *sarama.Config
// }
