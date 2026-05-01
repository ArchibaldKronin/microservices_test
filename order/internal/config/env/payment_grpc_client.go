package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentGRPCClientEnvConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST,required"`
	Port string `env:"PAYMENT_GRPC_PORT,required"`
}

type paymentGRPCClientConfig struct {
	config paymentGRPCClientEnvConfig
}

func NewPaymentGRPCClientConfig() (*paymentGRPCClientConfig, error) {
	var config paymentGRPCClientEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &paymentGRPCClientConfig{
		config: config,
	}, nil
}

func (cfg *paymentGRPCClientConfig) Address() string {
	return net.JoinHostPort(cfg.config.Host, cfg.config.Port)
}
