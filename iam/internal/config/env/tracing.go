package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type tracingEnvConfig struct {
	ServiceName        string `env:"OTEL_SERVICE_NAME"`
	ServiceEnvironment string `env:"OTEL_ENVIRONMENT"`
	ServiceVersion     string `env:"OTEL_SERVICE_VERSION"`
	OTLPHost           string `env:"OTLP_HOST,required"`
	OTLPPort           string `env:"OTLP_PORT,required"`
}

type tracingConfig struct {
	config tracingEnvConfig
}

func NewTracingConfig() (*tracingConfig, error) {
	var config tracingEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &tracingConfig{config: config}, nil
}

func (cfg *tracingConfig) CollectorEndpoint() string {
	return net.JoinHostPort(cfg.config.OTLPHost, cfg.config.OTLPPort)
}

func (cfg *tracingConfig) ServiceName() string {
	return cfg.config.ServiceName
}

func (cfg *tracingConfig) Environment() string {
	return cfg.config.ServiceEnvironment
}

func (cfg *tracingConfig) ServiceVersion() string {
	return cfg.config.ServiceVersion
}
