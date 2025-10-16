package usecase

import (
	"context"
	"errors"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

var ErrCategoryNotFound = errors.New("category not found")

type CategoryService struct {
	categoryRepo ports.CategoryRepository
}

func NewCategoryService(categoryRepo ports.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) GetAll(ctx context.Context) ([]domain.Category, error) {
	return s.categoryRepo.GetAll(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *CategoryService) Create(ctx context.Context, category *domain.Category) error {
	return s.categoryRepo.Create(ctx, category)
}

func (s *CategoryService) Update(ctx context.Context, category *domain.Category) error {
	existing, err := s.categoryRepo.GetByID(ctx, category.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCategoryNotFound
	}
	return s.categoryRepo.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCategoryNotFound
	}
	return s.categoryRepo.Delete(ctx, id)
}
