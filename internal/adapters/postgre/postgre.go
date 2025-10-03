package postgre

import (
	"database/sql"

	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

type Config struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	UserName string `env:"DB_USER" env-default:"admin"`
	Password string `env:"DB_PASSWORD" env-default:"admin123"`
	DBName   string `env:"DB_NAME" env-default:"restaurant_db"`
}

type Pool struct {
	DB *sql.DB
}

func NewPoolDB(db *sql.DB) ports.Postgres {
	return &Pool{
		DB: db,
	}
}
