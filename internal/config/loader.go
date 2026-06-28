package config

import (
	"github.com/caarlos0/env/v11"
)

var cfg Config

func LoadConfig(path string) (*Config, error) {
	// YAML CONFIGURATION

	// file, err := os.ReadFile(path)
	// if err != nil {
	// 	return nil, err
	// }
	// err = yaml.Unmarshal(file, &cfg)
	// if err != nil {
		// 	return nil, err
		// }
		
	// log.Println("Config: ", cfg)

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}
