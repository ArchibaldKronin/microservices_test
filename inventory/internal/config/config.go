package config

import (
	"os"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	Mongo         MongoConfig
	InventoryGRPC InventoryGRPCConfig

	IamConfig IamGRPCClientConfig

	Tracing TracingConfig
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

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	iamCfg, err := env.NewIamGRPCClientConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		Mongo:         mongoCfg,
		InventoryGRPC: inventoryGRPCCfg,
		IamConfig:     iamCfg,
		Tracing:       tracingCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
