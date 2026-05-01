package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type paymentGRPCConfig struct {
	config paymentGRPCEnvConfig
}

func NewPaymentGRPCConfig() (*paymentGRPCConfig, error) {
	var config paymentGRPCEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &paymentGRPCConfig{config: config}, nil
}

func (cfg *paymentGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.config.Host, cfg.config.Port)
}
