// file: internal/config/config.go

package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// PostgreConfig хранит конфигурацию для подключения к PostgreSQL.
type PostgreConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	UserName string `env:"DB_USER" env-default:"admin"`
	Password string `env:"DB_PASSWORD" env-default:"admin123"`
	DBName   string `env:"DB_NAME" env-default:"restaurant_db"`
}

// HTTPConfig хранит конфигурацию для HTTP-сервера.
type HTTPConfig struct {
	Port string `env:"HTTP_PORT" env-default:"8080"`
}

// Config - это главная структура конфигурации всего приложения.
type Config struct {
	Postgre PostgreConfig
	HTTP    HTTPConfig
}

// InitConfig читает конфигурацию из переменных окружения.
func InitConfig() (*Config, error) {
	var cfg Config

	// cleanenv автоматически найдет переменные окружения
	// (DB_HOST, DB_PORT, HTTP_PORT и т.д.) и заполнит структуру.
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
