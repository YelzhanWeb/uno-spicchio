package postgre

import "github.com/YelzhanWeb/uno-spicchio/internal/domain"

// AddIngredientToDish - добавить ингредиент к блюду
func (r *Pool) AddIngredientToDish(dishID, ingredientID int, qtyPerDish float64) error {
	stmt := `INSERT INTO dish_ingredients (dish_id, ingredient_id, qty_per_dish)
	         VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(stmt, dishID, ingredientID, qtyPerDish)
	return err
}

// GetIngredientsByDish - получить все ингредиенты для конкретного блюда
func (r *Pool) GetIngredientsByDish(dishID int) ([]domain.DishIngredient, error) {
	stmt := `SELECT dish_id, ingredient_id, qty_per_dish 
	         FROM dish_ingredients WHERE dish_id=$1`
	rows, err := r.DB.Query(stmt, dishID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []domain.DishIngredient
	for rows.Next() {
		var di domain.DishIngredient
		if err := rows.Scan(&di.DishID, &di.IngredientID, &di.QtyPerDish); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, di)
	}
	return ingredients, nil
}

// GetDishesByIngredient - получить все блюда, где используется ингредиент
func (r *Pool) GetDishesByIngredient(ingredientID int) ([]domain.DishIngredient, error) {
	stmt := `SELECT dish_id, ingredient_id, qty_per_dish 
	         FROM dish_ingredients WHERE ingredient_id=$1`
	rows, err := r.DB.Query(stmt, ingredientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes []domain.DishIngredient
	for rows.Next() {
		var di domain.DishIngredient
		if err := rows.Scan(&di.DishID, &di.IngredientID, &di.QtyPerDish); err != nil {
			return nil, err
		}
		dishes = append(dishes, di)
	}
	return dishes, nil
}

// UpdateDishIngredient - обновить количество ингредиента в блюде
func (r *Pool) UpdateDishIngredient(dishID, ingredientID int, qtyPerDish float64) error {
	stmt := `UPDATE dish_ingredients 
	         SET qty_per_dish=$1 
	         WHERE dish_id=$2 AND ingredient_id=$3`
	_, err := r.DB.Exec(stmt, qtyPerDish, dishID, ingredientID)
	return err
}

// RemoveIngredientFromDish - удалить ингредиент из блюда
func (r *Pool) RemoveIngredientFromDish(dishID, ingredientID int) error {
	stmt := `DELETE FROM dish_ingredients 
	         WHERE dish_id=$1 AND ingredient_id=$2`
	_, err := r.DB.Exec(stmt, dishID, ingredientID)
	return err
}
