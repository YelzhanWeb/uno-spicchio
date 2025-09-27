package postgre

import (
	"database/sql"

	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

type Pool struct {
	DB *sql.DB
}

func NewPoolDB(db *sql.DB) ports.Postgres {
	return &Pool{
		DB: db,
	}
}
