package config

import (
	"os"

	"github.com/ArchibaldKronin/microservices_test/assembly/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger LoggerConfig

	Kafka                 KafkaConfig
	ShipAssembledProducer ShipAssembledProducerConfig
	OrderPaidConsumer     OrderPaidConsumerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	shipAssembledProducerCfg, err := env.NewshipAssembledProducerConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                loggerCfg,
		Kafka:                 kafkaCfg,
		ShipAssembledProducer: shipAssembledProducerCfg,
		OrderPaidConsumer:     orderPaidConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
