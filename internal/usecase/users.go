// file: internal/usecase/users.go

package usecase

import (
	"context" // <-- Добавляем context
	"errors"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports" // <-- Добавляем ports
)

// UserService реализует бизнес-логику для пользователей.
type UserService struct {
	repo ports.UserRepository // <-- Зависимость от UserRepository, а не от всего пула
}

// NewUserService - конструктор для нашего сервиса пользователей.
func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Все методы теперь привязаны к (s *UserService), а не (s *Service)
// и принимают context.Context

func (s *UserService) CreateUser(ctx context.Context, username, password, role, photoKey string) (int, error) {
	if username == "" || password == "" {
		return 0, errors.New("username and password cannot be empty")
	}
	if role == "" {
		role = "waiter"
	}

	// TODO: захешировать пароль (bcrypt)
	passwordHash := password

	// Вызываем метод репозитория, а не пула: s.repo.CreateUser
	return s.repo.CreateUser(ctx, username, passwordHash, role, photoKey)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.PhotoKey != "" {
		// TODO: заменить на minioClient.PresignedGetObject()
		user.PhotoURL = "https://fake-minio.local/" + user.PhotoKey + "?expires=3600"
	}

	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	for i := range users {
		if users[i].PhotoKey != "" {
			users[i].PhotoURL = "https://fake-minio.local/" + users[i].PhotoKey + "?expires=3600"
		}
	}

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, username, password, role, photoKey string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if role == "" {
		role = "waiter"
	}

	passwordHash := password
	return s.repo.UpdateUser(ctx, id, username, passwordHash, role, photoKey)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}
