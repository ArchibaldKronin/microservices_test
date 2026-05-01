package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type inventoryGRPCConfig struct {
	config inventoryGRPCEnvConfig
}

func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	var config inventoryGRPCEnvConfig
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &inventoryGRPCConfig{config: config}, nil
}

func (cfg *inventoryGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.config.Host, cfg.config.Port)
}
