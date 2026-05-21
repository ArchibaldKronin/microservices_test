package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type metricsEnvConfig struct {
	CollectorInterval time.Duration `env:"OTEL_COLLECTOR_INTERVAL,required"`
	OTLPHost          string        `env:"OTLP_HOST,required"`
	OTLPPort          string        `env:"OTLP_PORT,required"`
	ServiceName       string        `env:"OTEL_SERVICE_NAME,required"`
}

type metricsConfig struct {
	config metricsEnvConfig
}

func NewMetricsConfig() (*metricsConfig, error) {
	var config metricsEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &metricsConfig{config: config}, nil
}

func (cfg *metricsConfig) CollectorEndpoint() string {
	return net.JoinHostPort(cfg.config.OTLPHost, cfg.config.OTLPPort)
}

func (cfg *metricsConfig) ServiceName() string {
	return cfg.config.ServiceName
}

func (cfg *metricsConfig) CollectorInterval() time.Duration {
	return cfg.config.CollectorInterval
}
