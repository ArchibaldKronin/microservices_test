package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type orderHTTPEnvConfig struct {
	Host        string        `env:"HTTP_HOST,required"`
	Port        string        `env:"HTTP_PORT,required"`
	ReadTimeout time.Duration `env:"HTTP_READ_TIMEOUT,required"`
}

type orderHTTPConfig struct {
	config orderHTTPEnvConfig
}

func NewOrderHTTPConfig() (*orderHTTPConfig, error) {
	var config orderHTTPEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &orderHTTPConfig{config: config}, nil
}

func (cfg *orderHTTPConfig) Address() string {
	return net.JoinHostPort(cfg.config.Host, cfg.config.Port)
}

func (cfg *orderHTTPConfig) Port() string {
	return cfg.config.Port
}

func (cfg *orderHTTPConfig) ReadTime() time.Duration {
	return cfg.config.ReadTimeout
}
