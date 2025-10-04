package usecase

import (
	"errors"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

func (s *Service) CreateUser(username, password, role, photoKey string) (int, error) {
	if username == "" || password == "" {
		return 0, errors.New("username and password cannot be empty")
	}
	if role == "" {
		role = "waiter"
	}

	// TODO: захешировать пароль (bcrypt)
	passwordHash := password

	return s.pool.CreateUser(username, passwordHash, role, photoKey)
}

func (s *Service) GetUserByID(id int) (*domain.User, error) {
	user, err := s.pool.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user.PhotoKey != "" {
		// TODO: заменить на minioClient.PresignedGetObject()
		user.PhotoURL = "https://fake-minio.local/" + user.PhotoKey + "?expires=3600"
	}

	return user, nil
}

func (s *Service) GetAllUsers() ([]domain.User, error) {
	users, err := s.pool.GetAllUsers()
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

func (s *Service) UpdateUser(id int, username, password, role, photoKey string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if role == "" {
		role = "waiter"
	}

	passwordHash := password
	return s.pool.UpdateUser(id, username, passwordHash, role, photoKey)
}

func (s *Service) DeleteUser(id int) error {
	return s.pool.DeleteUser(id)
}
