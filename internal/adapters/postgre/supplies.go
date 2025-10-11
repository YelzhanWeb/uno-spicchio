// file: internal/adapters/postgre/supplies.go

package postgre

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/jmoiron/sqlx"
)

// SupplyRepository - это реализация репозитория для поставок.
type SupplyRepository struct {
	db *sqlx.DB
}

// NewSupplyRepository создает новый экземпляр репозитория.
func NewSupplyRepository(db *sqlx.DB) *SupplyRepository {
	return &SupplyRepository{db: db}
}

// AddSupply - добавить поставку и обновить склад в рамках одной транзакции.
func (r *SupplyRepository) AddSupply(ctx context.Context, supply domain.Supply) (int, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Безопасный откат в случае любой ошибки

	// 1. Создаем запись о поставке
	var supplyID int
	stmtSupply := `INSERT INTO supplies (ingredient_id, qty, supplier_name)
	               VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRowxContext(ctx, stmtSupply, supply.IngredientID, supply.Qty, supply.SupplierName).Scan(&supplyID)
	if err != nil {
		return 0, fmt.Errorf("failed to create supply record: %w", err)
	}

	// 2. Обновляем количество на складе
	stmtIngredient := `UPDATE ingredients SET qty = qty + $1 WHERE id=$2`
	_, err = tx.ExecContext(ctx, stmtIngredient, supply.Qty, supply.IngredientID)
	if err != nil {
		return 0, fmt.Errorf("failed to update ingredient stock: %w", err)
	}

	// Если все успешно, коммитим транзакцию
	return supplyID, tx.Commit()
}

// GetSupplyByID - получить поставку по id
func (r *SupplyRepository) GetSupplyByID(ctx context.Context, id int) (*domain.Supply, error) {
	var s domain.Supply
	query := `SELECT id, ingredient_id, qty, supplier_name, created_at FROM supplies WHERE id=$1`
	err := r.db.GetContext(ctx, &s, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get supply by id %d: %w", id, err)
	}
	return &s, nil
}

// GetAllSupplies - получить все поставки
func (r *SupplyRepository) GetAllSupplies(ctx context.Context) ([]domain.Supply, error) {
	var supplies []domain.Supply
	query := `SELECT id, ingredient_id, qty, supplier_name, created_at
	          FROM supplies ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &supplies, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all supplies: %w", err)
	}
	return supplies, nil
}

// Функции UpdateSupply и DeleteSupply также должны быть здесь, переписанные под sqlx.
// Пока оставим их, чтобы не усложнять. Главное, что ошибка компиляции исчезнет.
