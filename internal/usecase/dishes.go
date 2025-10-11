// Файл: internal/usecase/dishes.go (НОВЫЙ ФАЙЛ)

package usecase

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

// DishService реализует бизнес-логику для блюд.
type DishService struct {
	repo ports.DishRepository
}

// NewDishService - конструктор для нашего сервиса блюд.
func NewDishService(repo ports.DishRepository) *DishService {
	return &DishService{repo: repo}
}

func (s *DishService) CreateDish(ctx context.Context, dish domain.Dish) (int, error) {
	// Здесь можно добавить валидацию, например:
	if dish.Name == "" {
		return 0, fmt.Errorf("dish name cannot be empty")
	}
	if dish.Price < 0 {
		return 0, fmt.Errorf("price cannot be negative")
	}
	// Устанавливаем значения по умолчанию
	dish.IsActive = true

	return s.repo.Create(ctx, dish)
}
