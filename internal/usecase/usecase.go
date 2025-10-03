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

func (s *Service) CreateUser(userName string, password string, role string) (int, error) {
	// TO DO Validate
	return s.pool.CreateUser(userName, password, role)
}
