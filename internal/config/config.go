package config

import (
	"github.com/YelzhanWeb/uno-spicchio/internal/adapters/postgre"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgre postgre.Config
}

func InitConfig() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	// Validation если нужна
	return &cfg, nil
}
