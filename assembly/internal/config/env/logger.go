package env

import (
	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level  string `env:"LOGGER_LEVEL,required"`
	AsJSON bool   `env:"LOGGER_AS_JSON,required"`
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
