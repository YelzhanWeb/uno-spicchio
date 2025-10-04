package postgre

import (
	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

// AddSupply - добавить поставку и обновить склад
func (r *Pool) AddSupply(ingredientID int, qty float64, supplierName string) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	stmtSupply := `INSERT INTO supplies (ingredient_id, qty, supplier_name)
                   VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(stmtSupply, ingredientID, qty, supplierName).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	stmtIngredient := `UPDATE ingredients SET qty = qty + $1 WHERE id=$2`
	_, err = tx.Exec(stmtIngredient, qty, ingredientID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

// UpdateSupply - обновить поставку и откорректировать склад
func (r *Pool) UpdateSupply(id int, newQty float64, newSupplier string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// 1. получаем старое значение
	var oldQty float64
	var ingredientID int
	stmtOld := `SELECT ingredient_id, qty FROM supplies WHERE id=$1`
	err = tx.QueryRow(stmtOld, id).Scan(&ingredientID, &oldQty)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. обновляем поставку
	stmtUpdate := `UPDATE supplies SET qty=$1, supplier_name=$2 WHERE id=$3`
	_, err = tx.Exec(stmtUpdate, newQty, newSupplier, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3. корректируем склад
	diff := newQty - oldQty
	stmtIngredient := `UPDATE ingredients SET qty = qty + $1 WHERE id=$2`
	_, err = tx.Exec(stmtIngredient, diff, ingredientID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// DeleteSupply - удалить поставку и откорректировать склад
func (r *Pool) DeleteSupply(id int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// 1. получаем данные поставки
	var qty float64
	var ingredientID int
	stmtOld := `SELECT ingredient_id, qty FROM supplies WHERE id=$1`
	err = tx.QueryRow(stmtOld, id).Scan(&ingredientID, &qty)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. удаляем поставку
	stmtDelete := `DELETE FROM supplies WHERE id=$1`
	_, err = tx.Exec(stmtDelete, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3. уменьшаем количество на складе
	stmtIngredient := `UPDATE ingredients SET qty = qty - $1 WHERE id=$2`
	_, err = tx.Exec(stmtIngredient, qty, ingredientID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetSupplyByID - получить поставку по id
func (r *Pool) GetSupplyByID(id int) (*domain.Supply, error) {
	stmt := `SELECT id, ingredient_id, qty, supplier_name, created_at FROM supplies WHERE id=$1`
	row := r.DB.QueryRow(stmt, id)

	var s domain.Supply
	err := row.Scan(&s.ID, &s.IngredientID, &s.Qty, &s.SupplierName, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetAllSupplies - получить все поставки
func (r *Pool) GetAllSupplies() ([]domain.Supply, error) {
	stmt := `SELECT id, ingredient_id, qty, supplier_name, created_at 
             FROM supplies ORDER BY created_at DESC`
	rows, err := r.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var supplies []domain.Supply
	for rows.Next() {
		var s domain.Supply
		if err := rows.Scan(&s.ID, &s.IngredientID, &s.Qty, &s.SupplierName, &s.CreatedAt); err != nil {
			return nil, err
		}
		supplies = append(supplies, s)
	}
	return supplies, nil
}

// GetSuppliesByIngredient - история поставок по ингредиенту
func (r *Pool) GetSuppliesByIngredient(ingredientID int) ([]domain.Supply, error) {
	stmt := `SELECT id, ingredient_id, qty, supplier_name, created_at 
             FROM supplies WHERE ingredient_id=$1 ORDER BY created_at DESC`
	rows, err := r.DB.Query(stmt, ingredientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var supplies []domain.Supply
	for rows.Next() {
		var s domain.Supply
		if err := rows.Scan(&s.ID, &s.IngredientID, &s.Qty, &s.SupplierName, &s.CreatedAt); err != nil {
			return nil, err
		}
		supplies = append(supplies, s)
	}
	return supplies, nil
}
