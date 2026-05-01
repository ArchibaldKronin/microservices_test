package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	Host     string `env:"MONGO_HOST,required"`
	Port     string `env:"MONGO_PORT,required"`
	Database string `env:"MONGO_DATABASE,required"`
	AuthDB   string `env:"MONGO_AUTH_DB,required"`
	Username string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
}

type mongoConfig struct {
	config mongoEnvConfig
}

func NewMongoConfig() (*mongoConfig, error) {
	var config mongoEnvConfig
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &mongoConfig{config: config}, nil
}

func (cfg *mongoConfig) URI() string {
	fmt.Print(cfg.config.Host)
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.config.Username,
		cfg.config.Password,
		cfg.config.Host,
		cfg.config.Port,
		cfg.config.Database,
		cfg.config.AuthDB,
	)
}

func (cfg *mongoConfig) DatabaseName() string {
	return cfg.config.Database
}
