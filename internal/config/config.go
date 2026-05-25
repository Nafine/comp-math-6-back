package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host string `env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port string `env:"HTTP_PORT" env-default:"8080"`
}

func Get() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
