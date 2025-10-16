package postgre

import (
	"context"
	"database/sql"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

type SupplyRepository struct {
	db *sql.DB
}

func NewSupplyRepository(db *sql.DB) *SupplyRepository {
	return &SupplyRepository{db: db}
}

func (r *SupplyRepository) Create(ctx context.Context, supply *domain.Supply) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert supply
	query := `
		INSERT INTO supplies (ingredient_id, qty, supplier_name)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err = tx.QueryRowContext(ctx, query,
		supply.IngredientID, supply.Qty, supply.SupplierName,
	).Scan(&supply.ID, &supply.CreatedAt)
	if err != nil {
		return err
	}

	// Update ingredient quantity
	updateQuery := `UPDATE ingredients SET qty = qty + $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, updateQuery, supply.Qty, supply.IngredientID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SupplyRepository) GetAll(ctx context.Context) ([]domain.Supply, error) {
	query := `
		SELECT s.id, s.ingredient_id, s.qty, s.supplier_name, s.created_at,
		       i.name, i.unit
		FROM supplies s
		JOIN ingredients i ON s.ingredient_id = i.id
		ORDER BY s.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var supplies []domain.Supply
	for rows.Next() {
		var supply domain.Supply
		supply.Ingredient = &domain.Ingredient{}

		if err := rows.Scan(
			&supply.ID, &supply.IngredientID, &supply.Qty, &supply.SupplierName, &supply.CreatedAt,
			&supply.Ingredient.Name, &supply.Ingredient.Unit,
		); err != nil {
			return nil, err
		}
		supplies = append(supplies, supply)
	}

	return supplies, rows.Err()
}

func (r *SupplyRepository) GetByID(ctx context.Context, id int) (*domain.Supply, error) {
	query := `
		SELECT s.id, s.ingredient_id, s.qty, s.supplier_name, s.created_at,
		       i.name, i.unit
		FROM supplies s
		JOIN ingredients i ON s.ingredient_id = i.id
		WHERE s.id = $1`

	supply := &domain.Supply{Ingredient: &domain.Ingredient{}}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&supply.ID, &supply.IngredientID, &supply.Qty, &supply.SupplierName, &supply.CreatedAt,
		&supply.Ingredient.Name, &supply.Ingredient.Unit,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return supply, err
}

func (r *SupplyRepository) GetByIngredientID(ctx context.Context, ingredientID int) ([]domain.Supply, error) {
	query := `
		SELECT id, ingredient_id, qty, supplier_name, created_at
		FROM supplies
		WHERE ingredient_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, ingredientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var supplies []domain.Supply
	for rows.Next() {
		var supply domain.Supply
		if err := rows.Scan(
			&supply.ID, &supply.IngredientID, &supply.Qty, &supply.SupplierName, &supply.CreatedAt,
		); err != nil {
			return nil, err
		}
		supplies = append(supplies, supply)
	}

	return supplies, rows.Err()
}
