package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamGRPCClientEnvConfig struct {
	Host string `env:"IAM_GRPC_HOST,required"`
	Port string `env:"IAM_GRPC_PORT,required"`
}

type iamGRPCClientConfig struct {
	config iamGRPCClientEnvConfig
}

func NewIamGRPCClientConfig() (*iamGRPCClientConfig, error) {
	var config iamGRPCClientEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &iamGRPCClientConfig{
		config: config,
	}, nil
}

func (cfg *iamGRPCClientConfig) Address() string {
	return net.JoinHostPort(cfg.config.Host, cfg.config.Port)
}
