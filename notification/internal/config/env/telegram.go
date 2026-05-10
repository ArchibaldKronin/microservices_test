package env

import "github.com/caarlos0/env/v11"

type telegramEnvConfig struct {
	Token string `env:"TELEGRAM_BOT_TOKEN,required"`
}

type telegramConfig struct {
	config telegramEnvConfig
}

func NewTelegramConfig() (*telegramConfig, error) {
	var config telegramEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &telegramConfig{config: config}, nil
}

func (cfg *telegramConfig) Token() string {
	return cfg.config.Token
}
