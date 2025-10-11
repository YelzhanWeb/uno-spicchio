// file: internal/usecase/categories.go

package usecase

import (
	"context" // <-- Добавляем context
	"errors"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports" // <-- Добавляем ports
)

// CategoryService реализует бизнес-логику для категорий.
type CategoryService struct {
	repo ports.CategoryRepository
}

// NewCategoryService - конструктор для нашего сервиса категорий.
func NewCategoryService(repo ports.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// Все методы теперь привязаны к (s *CategoryService) и принимают context

func (s *CategoryService) CreateCategory(ctx context.Context, name string) (int, error) {
	if name == "" {
		return 0, errors.New("category name cannot be empty")
	}
	// Вызываем метод репозитория: s.repo.CreateCategory
	return s.repo.CreateCategory(ctx, name)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.GetAllCategories(ctx)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id int, name string) error {
	if name == "" {
		return errors.New("category name cannot be empty")
	}
	return s.repo.UpdateCategory(ctx, id, name)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.DeleteCategory(ctx, id)
}
