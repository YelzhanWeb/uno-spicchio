package usecase

import (
	"context"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

type SupplyService struct {
	supplyRepo ports.SupplyRepository
}

func NewSupplyService(supplyRepo ports.SupplyRepository) *SupplyService {
	return &SupplyService{supplyRepo: supplyRepo}
}

func (s *SupplyService) Create(ctx context.Context, supply *domain.Supply) error {
	return s.supplyRepo.Create(ctx, supply)
}

func (s *SupplyService) GetAll(ctx context.Context) ([]domain.Supply, error) {
	return s.supplyRepo.GetAll(ctx)
}

func (s *SupplyService) GetByIngredientID(ctx context.Context, ingredientID int) ([]domain.Supply, error) {
	return s.supplyRepo.GetByIngredientID(ctx, ingredientID)
}
