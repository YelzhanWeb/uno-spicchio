package usecase

import (
	"context"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

type IngredientService struct {
	ingredientRepo ports.IngredientRepository
}

func NewIngredientService(ingredientRepo ports.IngredientRepository) *IngredientService {
	return &IngredientService{ingredientRepo: ingredientRepo}
}

func (s *IngredientService) GetAll(ctx context.Context) ([]domain.Ingredient, error) {
	return s.ingredientRepo.GetAll(ctx)
}

func (s *IngredientService) GetByID(ctx context.Context, id int) (*domain.Ingredient, error) {
	return s.ingredientRepo.GetByID(ctx, id)
}

func (s *IngredientService) GetLowStock(ctx context.Context) ([]domain.Ingredient, error) {
	return s.ingredientRepo.GetLowStock(ctx)
}

func (s *IngredientService) Create(ctx context.Context, ingredient *domain.Ingredient) error {
	return s.ingredientRepo.Create(ctx, ingredient)
}

func (s *IngredientService) Update(ctx context.Context, ingredient *domain.Ingredient) error {
	return s.ingredientRepo.Update(ctx, ingredient)
}

func (s *IngredientService) Delete(ctx context.Context, id int) error {
	return s.ingredientRepo.Delete(ctx, id)
}
