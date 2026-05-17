package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type gRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type gRPCConfig struct {
	config gRPCEnvConfig
}

func NewGRPCConfig() (*gRPCConfig, error) {
	var config gRPCEnvConfig
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &gRPCConfig{config: config}, nil
}

func (cfg *gRPCConfig) Address() string {
	return net.JoinHostPort(cfg.config.Host, cfg.config.Port)
}
