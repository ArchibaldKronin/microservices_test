package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level              string `env:"LOGGER_LEVEL,required"`
	AsJSON             bool   `env:"LOGGER_AS_JSON,required"`
	ServiceName        string `env:"LOGGER_SERVICE_NAME,required"`
	ServiceEnvironment string `env:"LOGGER_SERVICE_ENVIRONMENT,required"`
	EnableOTLP         bool   `env:"LOGGER_ENABLE_OTLP,required"`
	OTLPHost           string `env:"OTLP_HOST,required"`
	OTLPPort           string `env:"OTLP_PORT,required"`
}

type loggerConfig struct {
	config loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var config loggerEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &loggerConfig{config: config}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.config.Level
}

func (cfg *loggerConfig) AsJSON() bool {
	return cfg.config.AsJSON
}

func (cfg *loggerConfig) OTLPAddress() string {
	return net.JoinHostPort(cfg.config.OTLPHost, cfg.config.OTLPPort)
}

func (cfg *loggerConfig) ServiceName() string {
	return cfg.config.ServiceName
}

func (cfg *loggerConfig) ServiceEnvironment() string {
	return cfg.config.ServiceEnvironment
}

func (cfg *loggerConfig) EnableOTLP() bool {
	return cfg.config.EnableOTLP
}
