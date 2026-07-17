package config

import (
	"github.com/caarlos0/env/v11"
)

var cfg Config

func LoadEnv() (*Config, error) {
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}
