package postgre

import "github.com/YelzhanWeb/uno-spicchio/internal/domain"

// CreateIngredient - добавление нового ингредиента
func (r *Pool) CreateIngredient(name, unit string, qty, minQty float64) (int, error) {
	stmt := `INSERT INTO ingredients (name, unit, qty, min_qty)
	         VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.DB.QueryRow(stmt, name, unit, qty, minQty).Scan(&id)
	return id, err
}

// GetIngredientByID - получить ингредиент по id
func (r *Pool) GetIngredientByID(id int) (*domain.Ingredient, error) {
	stmt := `SELECT id, name, unit, qty, min_qty FROM ingredients WHERE id = $1`
	row := r.DB.QueryRow(stmt, id)

	var ingredient domain.Ingredient
	err := row.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Qty, &ingredient.MinQty)
	if err != nil {
		return nil, err
	}
	return &ingredient, nil
}

// GetAllIngredients - получить все ингредиенты
func (r *Pool) GetAllIngredients() ([]domain.Ingredient, error) {
	stmt := `SELECT id, name, unit, qty, min_qty FROM ingredients`
	rows, err := r.DB.Query(stmt)
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
	return ingredients, nil
}

// UpdateIngredient - обновить ингредиент
func (r *Pool) UpdateIngredient(id int, name, unit string, qty, minQty float64) error {
	stmt := `UPDATE ingredients 
	         SET name=$1, unit=$2, qty=$3, min_qty=$4 
	         WHERE id=$5`
	_, err := r.DB.Exec(stmt, name, unit, qty, minQty, id)
	return err
}

// DeleteIngredient - удалить ингредиент по id
func (r *Pool) DeleteIngredient(id int) error {
	stmt := `DELETE FROM ingredients WHERE id=$1`
	_, err := r.DB.Exec(stmt, id)
	return err
}
