package usecase

import (
	"context"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/pkg/hash"
	"github.com/YelzhanWeb/uno-spicchio/pkg/jwt"
)

type AuthService struct {
	userRepo     ports.UserRepository
	tokenManager *jwt.TokenManager
}

func NewAuthService(userRepo ports.UserRepository, tokenManager *jwt.TokenManager) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}

	if user == nil || !hash.Verify(password, user.PasswordHash) {
		return "", nil, domain.ErrInvalidCredentials
	}
	if !user.IsActive {
		return "", nil, domain.ErrUserNotActive
	}

	token, err := s.tokenManager.Generate(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID int) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
