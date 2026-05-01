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

	OrderHTTP OrderHTTPConfigConfig

	Postgres PostgresConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	LoggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	InventoryClientCfg, err := env.NewInventoryGRPCClientConfig()
	if err != nil {
		return err
	}

	PaymentClientCfg, err := env.NewPaymentGRPCClientConfig()
	if err != nil {
		return err
	}

	OrderConfig, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	PostgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:              LoggerCfg,
		InventoryGRPCClient: InventoryClientCfg,
		PaymentGRPCClient:   PaymentClientCfg,
		OrderHTTP:           OrderConfig,
		Postgres:            PostgresCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
