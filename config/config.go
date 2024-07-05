package config

import (
	"github.com/traefik/paerser/file"
)

type Config struct {
	ListenAddress      string
	Telegram           TelegramConfig
	TemplatesDirectory string
}

type TelegramConfig struct {
	Token     string
	Receiver  string
	ParseMode string
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	err := file.Decode(path, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
