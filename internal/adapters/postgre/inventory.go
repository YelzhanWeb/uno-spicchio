package postgre

import (
	"context"
	"database/sql"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

type IngredientRepository struct {
	db *sql.DB
}

func NewIngredientRepository(db *sql.DB) *IngredientRepository {
	return &IngredientRepository{db: db}
}

func (r *IngredientRepository) GetAll(ctx context.Context) ([]domain.Ingredient, error) {
	query := `SELECT id, name, unit, qty, min_qty FROM ingredients ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []domain.Ingredient
	for rows.Next() {
		var ing domain.Ingredient
		if err := rows.Scan(&ing.ID, &ing.Name, &ing.Unit, &ing.Qty, &ing.MinQty); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ing)
	}

	return ingredients, rows.Err()
}

func (r *IngredientRepository) GetByID(ctx context.Context, id int) (*domain.Ingredient, error) {
	query := `SELECT id, name, unit, qty, min_qty FROM ingredients WHERE id = $1`

	ing := &domain.Ingredient{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&ing.ID, &ing.Name, &ing.Unit, &ing.Qty, &ing.MinQty)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return ing, err
}

func (r *IngredientRepository) GetLowStock(ctx context.Context) ([]domain.Ingredient, error) {
	query := `SELECT id, name, unit, qty, min_qty FROM ingredients WHERE qty <= min_qty ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []domain.Ingredient
	for rows.Next() {
		var ing domain.Ingredient
		if err := rows.Scan(&ing.ID, &ing.Name, &ing.Unit, &ing.Qty, &ing.MinQty); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ing)
	}

	return ingredients, rows.Err()
}

func (r *IngredientRepository) Create(ctx context.Context, ing *domain.Ingredient) error {
	query := `
		INSERT INTO ingredients (name, unit, qty, min_qty)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query, ing.Name, ing.Unit, ing.Qty, ing.MinQty).Scan(&ing.ID)
}

func (r *IngredientRepository) Update(ctx context.Context, ing *domain.Ingredient) error {
	query := `
		UPDATE ingredients 
		SET name = $1, unit = $2, qty = $3, min_qty = $4
		WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query, ing.Name, ing.Unit, ing.Qty, ing.MinQty, ing.ID)
	return err
}

func (r *IngredientRepository) UpdateQuantity(ctx context.Context, id int, qty float64) error {
	query := `UPDATE ingredients SET qty = qty + $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, qty, id)
	return err
}

func (r *IngredientRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM ingredients WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
