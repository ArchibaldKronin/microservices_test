package config

import (
	"os"

	"github.com/ArchibaldKronin/microservices_test/notification/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger LoggerConfig

	Kafka                  KafkaConfig
	OrderPaidConsumer      ConsumerConfig
	OrderAssembledConsumer ConsumerConfig
	Telegram               TelegramConfig
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

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	tgCfg, err := env.NewTelegramConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		Kafka:                  kafkaCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
		Telegram:               tgCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
