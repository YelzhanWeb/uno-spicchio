// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	MinIO    MinIOConfig
	JWT      JWTConfig
	Env      string
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type MinIOConfig struct {
	Endpoint     string
	AccessKey    string
	SecretKey    string
	UseSSL       bool
	BucketDishes string
	BucketUsers  string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "restaurant_user"),
			Password: getEnv("DB_PASSWORD", "restaurant_password"),
			DBName:   getEnv("DB_NAME", "restaurant_crm"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		MinIO: MinIOConfig{
			Endpoint:     getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:    getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:    getEnv("MINIO_SECRET_KEY", "minioadmin"),
			UseSSL:       getEnvBool("MINIO_USE_SSL", false),
			BucketDishes: getEnv("MINIO_BUCKET_DISHES", "dishes"),
			BucketUsers:  getEnv("MINIO_BUCKET_USERS", "users"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "change-me-in-production"),
			ExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		},
		Env: getEnv("ENV", "development"),
	}

	if cfg.JWT.Secret == "change-me-in-production" && cfg.Env == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return cfg, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func (c *JWTConfig) ExpirationDuration() time.Duration {
	return time.Duration(c.ExpirationHours) * time.Hour
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultVal
}
