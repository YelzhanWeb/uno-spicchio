// file: internal/adapters/postgre/dish_ingredients.go

package postgre

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/jmoiron/sqlx"
)

// DishIngredientRepository - это реализация репозитория для связи блюд и ингредиентов.
type DishIngredientRepository struct {
	db *sqlx.DB
}

// NewDishIngredientRepository создает новый экземпляр репозитория.
func NewDishIngredientRepository(db *sqlx.DB) *DishIngredientRepository {
	return &DishIngredientRepository{db: db}
}

// AddIngredientToDish - добавить ингредиент к блюду
func (r *DishIngredientRepository) AddIngredientToDish(ctx context.Context, dishID, ingredientID int, qtyPerDish float64) error {
	query := `INSERT INTO dish_ingredients (dish_id, ingredient_id, qty_per_dish)
	          VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, dishID, ingredientID, qtyPerDish)
	if err != nil {
		return fmt.Errorf("failed to add ingredient %d to dish %d: %w", ingredientID, dishID, err)
	}
	return nil
}

// GetIngredientsByDish - получить все ингредиенты для конкретного блюда
func (r *DishIngredientRepository) GetIngredientsByDish(ctx context.Context, dishID int) ([]domain.DishIngredient, error) {
	var ingredients []domain.DishIngredient
	query := `SELECT dish_id, ingredient_id, qty_per_dish
	          FROM dish_ingredients WHERE dish_id=$1`

	err := r.db.SelectContext(ctx, &ingredients, query, dishID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients for dish %d: %w", dishID, err)
	}
	return ingredients, nil
}

// GetDishesByIngredient - получить все блюда, где используется ингредиент
func (r *DishIngredientRepository) GetDishesByIngredient(ctx context.Context, ingredientID int) ([]domain.DishIngredient, error) {
	var dishes []domain.DishIngredient
	query := `SELECT dish_id, ingredient_id, qty_per_dish
	          FROM dish_ingredients WHERE ingredient_id=$1`

	err := r.db.SelectContext(ctx, &dishes, query, ingredientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dishes for ingredient %d: %w", ingredientID, err)
	}
	return dishes, nil
}

// UpdateDishIngredient - обновить количество ингредиента в блюде
func (r *DishIngredientRepository) UpdateDishIngredient(ctx context.Context, dishID, ingredientID int, qtyPerDish float64) error {
	query := `UPDATE dish_ingredients
	          SET qty_per_dish=$1
	          WHERE dish_id=$2 AND ingredient_id=$3`
	_, err := r.db.ExecContext(ctx, query, qtyPerDish, dishID, ingredientID)
	if err != nil {
		return fmt.Errorf("failed to update ingredient %d for dish %d: %w", ingredientID, dishID, err)
	}
	return nil
}

// RemoveIngredientFromDish - удалить ингредиент из блюда
func (r *DishIngredientRepository) RemoveIngredientFromDish(ctx context.Context, dishID, ingredientID int) error {
	query := `DELETE FROM dish_ingredients
	          WHERE dish_id=$1 AND ingredient_id=$2`
	_, err := r.db.ExecContext(ctx, query, dishID, ingredientID)
	if err != nil {
		return fmt.Errorf("failed to remove ingredient %d from dish %d: %w", ingredientID, dishID, err)
	}
	return nil
}
