package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryGRPCClientEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

type inventoryGRPCClientConfig struct {
	config inventoryGRPCClientEnvConfig
}

func NewInventoryGRPCClientConfig() (*inventoryGRPCClientConfig, error) {
	var config inventoryGRPCClientEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &inventoryGRPCClientConfig{
		config: config,
	}, nil
}

func (cfg *inventoryGRPCClientConfig) Address() string {
	return net.JoinHostPort(cfg.config.Host, cfg.config.Port)
}
