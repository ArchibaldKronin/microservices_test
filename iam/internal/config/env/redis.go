package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	Host              string        `env:"REDIS_HOST,required"`
	Port              string        `env:"REDIS_PORT,required"`
	ConnectionTimeout time.Duration `env:"REDIS_CONNECTION_TIMEOUT,required"`
	MaxIdle           int           `env:"REDIS_MAX_IDLE,required"`
	IdleTimeout       time.Duration `env:"REDIS_IDLE_TIMEOUT,required"`
	CacheTTL          time.Duration `env:"SESSION_TTL,required"`
}

type redisConfig struct {
	config redisEnvConfig
}

func NewRedisConfig() (*redisConfig, error) {
	var config redisEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &redisConfig{config: config}, nil
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.config.Host, cfg.config.Port)
}

func (cfg *redisConfig) ConnectionTimeout() time.Duration {
	return cfg.config.ConnectionTimeout
}

func (cfg *redisConfig) MaxIdle() int {
	return cfg.config.MaxIdle
}

func (cfg *redisConfig) IdleTimeout() time.Duration {
	return cfg.config.IdleTimeout
}

func (cfg *redisConfig) CacheTTL() time.Duration {
	return cfg.config.CacheTTL
}
