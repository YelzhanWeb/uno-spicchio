package ports

import "github.com/YelzhanWeb/uno-spicchio/internal/domain"

type Postgres interface {
	CreateUser(username, passwordHash, role, photoKey string) (int, error)
	GetUserByID(id int) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	UpdateUser(id int, username, passwordHash, role, photoKey string) error
	DeleteUser(id int) error
}
type Service interface {
	CreateUser(username, password, role, photoKey string) (int, error)
	GetUserByID(id int) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	UpdateUser(id int, username, password, role, photoKey string) error
	DeleteUser(id int) error
}

type HttpHandlers interface {
	PostUserHandler()
}
