// internal/core/services/dish_service.go
package usecase

import (
	"context"
	"errors"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

var ErrDishNotFound = errors.New("dish not found")

type DishService struct {
	dishRepo ports.DishRepository
}

func NewDishService(dishRepo ports.DishRepository) *DishService {
	return &DishService{dishRepo: dishRepo}
}

func (s *DishService) GetAll(ctx context.Context, activeOnly bool) ([]domain.Dish, error) {
	return s.dishRepo.GetAll(ctx, activeOnly)
}

func (s *DishService) GetByID(ctx context.Context, id int) (*domain.Dish, error) {
	dish, err := s.dishRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if dish == nil {
		return nil, ErrDishNotFound
	}
	return dish, nil
}

func (s *DishService) GetByCategoryID(ctx context.Context, categoryID int) ([]domain.Dish, error) {
	return s.dishRepo.GetByCategoryID(ctx, categoryID)
}

func (s *DishService) Create(ctx context.Context, dish *domain.Dish) error {
	dish.IsActive = true
	return s.dishRepo.Create(ctx, dish)
}

func (s *DishService) Update(ctx context.Context, dish *domain.Dish) error {
	existing, err := s.dishRepo.GetByID(ctx, dish.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrDishNotFound
	}

	return s.dishRepo.Update(ctx, dish)
}

func (s *DishService) Delete(ctx context.Context, id int) error {
	return s.dishRepo.Delete(ctx, id)
}

func (s *DishService) GetIngredients(ctx context.Context, dishID int) ([]domain.DishIngredient, error) {
	return s.dishRepo.GetIngredients(ctx, dishID)
}

func (s *DishService) AddIngredient(ctx context.Context, dishIngredient *domain.DishIngredient) error {
	return s.dishRepo.AddIngredient(ctx, dishIngredient)
}

func (s *DishService) RemoveIngredient(ctx context.Context, dishID, ingredientID int) error {
	return s.dishRepo.RemoveIngredient(ctx, dishID, ingredientID)
}
