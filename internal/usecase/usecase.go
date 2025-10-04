package usecase

import "github.com/YelzhanWeb/uno-spicchio/internal/ports"

type Service struct {
	pool ports.Postgres
}

func NewService(pool ports.Postgres) ports.Service {
	return &Service{
		pool: pool,
	}
}
