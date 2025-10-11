// file: internal/adapters/postgre/inventory.go

package postgre

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/jmoiron/sqlx"
)

// IngredientRepository - это реализация репозитория для ингредиентов.
type IngredientRepository struct {
	db *sqlx.DB
}

// NewIngredientRepository создает новый экземпляр репозитория.
func NewIngredientRepository(db *sqlx.DB) *IngredientRepository {
	return &IngredientRepository{db: db}
}

// CreateIngredient - добавление нового ингредиента
func (r *IngredientRepository) CreateIngredient(ctx context.Context, ingredient domain.Ingredient) (int, error) {
	query := `INSERT INTO ingredients (name, unit, qty, min_qty)
	          VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.db.QueryRowxContext(ctx, query, ingredient.Name, ingredient.Unit, ingredient.Qty, ingredient.MinQty).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create ingredient: %w", err)
	}
	return id, nil
}

// GetIngredientByID - получить ингредиент по id
func (r *IngredientRepository) GetIngredientByID(ctx context.Context, id int) (*domain.Ingredient, error) {
	var ingredient domain.Ingredient
	query := `SELECT id, name, unit, qty, min_qty FROM ingredients WHERE id = $1`

	err := r.db.GetContext(ctx, &ingredient, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredient by id %d: %w", id, err)
	}
	return &ingredient, nil
}

// GetAllIngredients - получить все ингредиенты
func (r *IngredientRepository) GetAllIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	var ingredients []domain.Ingredient
	query := `SELECT id, name, unit, qty, min_qty FROM ingredients`

	err := r.db.SelectContext(ctx, &ingredients, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all ingredients: %w", err)
	}
	return ingredients, nil
}

// UpdateIngredient - обновить ингредиент
func (r *IngredientRepository) UpdateIngredient(ctx context.Context, ingredient domain.Ingredient) error {
	query := `UPDATE ingredients
	          SET name=$1, unit=$2, qty=$3, min_qty=$4
	          WHERE id=$5`
	_, err := r.db.ExecContext(ctx, query, ingredient.Name, ingredient.Unit, ingredient.Qty, ingredient.MinQty, ingredient.ID)
	if err != nil {
		return fmt.Errorf("failed to update ingredient with id %d: %w", ingredient.ID, err)
	}
	return nil
}

// DeleteIngredient - удалить ингредиент по id
func (r *IngredientRepository) DeleteIngredient(ctx context.Context, id int) error {
	query := `DELETE FROM ingredients WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ingredient with id %d: %w", id, err)
	}
	return nil
}
