package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host          string `env:"POSTGRES_HOST,required"`
	Port          string `env:"POSTGRES_PORT,required"`
	Username      string `env:"POSTGRES_USER,required"`
	Password      string `env:"POSTGRES_PASSWORD,required"`
	Database      string `env:"POSTGRES_DB,required"`
	MigrationsDir string `env:"MIGRATION_DIRECTORY,required"`
	SSLMode       string `env:"POSTGRES_SSL_MODE,required"`
}

type postgresConfig struct {
	config postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var config postgresEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &postgresConfig{config: config}, nil
}

func (cfg *postgresConfig) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.config.Username,
		cfg.config.Password,
		cfg.config.Host,
		cfg.config.Port,
		cfg.config.Database,
		cfg.config.SSLMode,
	)
}

func (cfg *postgresConfig) MigrationsDir() string {
	return cfg.config.MigrationsDir
}
