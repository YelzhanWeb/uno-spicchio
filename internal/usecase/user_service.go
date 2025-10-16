package usecase

import (
	"context"
	"errors"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/pkg/hash"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Create(ctx context.Context, user *domain.User, password string) error {
	existing, err := s.userRepo.GetByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrUserExists
	}

	passwordHash, err := hash.Hash(password)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	user.IsActive = true

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, user *domain.User) error {
	existing, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Update(ctx, user)
}

func (s *UserService) UpdatePassword(ctx context.Context, userID int, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	passwordHash, err := hash.Hash(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(ctx, id)
}
