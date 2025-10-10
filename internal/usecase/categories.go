package usecase

import (
	"errors"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

func (s *Service) CreateCategory(name string) (int, error) {
	if name == "" {
		return 0, errors.New("category name cannot be empty")
	}
	return s.pool.CreateCategory(name)
}

func (s *Service) GetCategoryByID(id int) (*domain.Category, error) {
	return s.pool.GetCategoryByID(id)
}

func (s *Service) GetAllCategories() ([]domain.Category, error) {
	return s.pool.GetAllCategories()
}

func (s *Service) UpdateCategory(id int, name string) error {
	if name == "" {
		return errors.New("category name cannot be empty")
	}
	return s.pool.UpdateCategory(id, name)
}

func (s *Service) DeleteCategory(id int) error {
	return s.pool.DeleteCategory(id)
}
