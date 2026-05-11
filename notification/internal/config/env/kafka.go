package env

import "github.com/caarlos0/env/v11"

type kafkaEnvConfig struct {
	Brokers []string `env:"KAFKA_BROKERS,required"`
}

type kafkaConfig struct {
	config kafkaEnvConfig
}

func NewKafkaConfig() (*kafkaConfig, error) {
	var config kafkaEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &kafkaConfig{config: config}, nil
}

func (cfg *kafkaConfig) Brokers() []string {
	return cfg.config.Brokers
}
