package config

import (
	"os"

	"github.com/ArchibaldKronin/microservices_test/order/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger LoggerConfig

	InventoryGRPCClient InventoryGRPCClientConfig
	PaymentGRPCClient   PaymentGRPCClientConfig

	IamConfig IamGRPCClientConfig

	OrderHTTP OrderHTTPConfigConfig

	Postgres PostgresConfig

	Kafka                  KafkaConfig
	OrderPaidProducer      OrderPaidProducerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
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

	inventoryClientCfg, err := env.NewInventoryGRPCClientConfig()
	if err != nil {
		return err
	}

	paymentClientCfg, err := env.NewPaymentGRPCClientConfig()
	if err != nil {
		return err
	}

	iamClientCfg, err := env.NewIamGRPCClientConfig()
	if err != nil {
		return err
	}

	orderConfig, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}
	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		InventoryGRPCClient:    inventoryClientCfg,
		PaymentGRPCClient:      paymentClientCfg,
		OrderHTTP:              orderConfig,
		Postgres:               postgresCfg,
		Kafka:                  kafkaCfg,
		OrderPaidProducer:      orderPaidProducerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
		IamConfig:              iamClientCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
