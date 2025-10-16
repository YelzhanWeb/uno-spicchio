// internal/adapters/postgres/dish_repository.go
package postgre

import (
	"context"
	"database/sql"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

type DishRepository struct {
	db *sql.DB
}

func NewDishRepository(db *sql.DB) *DishRepository {
	return &DishRepository{db: db}
}

func (r *DishRepository) GetAll(ctx context.Context, activeOnly bool) ([]domain.Dish, error) {
	query := `
		SELECT d.id, d.category_id, d.name, d.description, d.price, d.photo_url, d.is_active,
		       c.id, c.name
		FROM dishes d
		LEFT JOIN categories c ON d.category_id = c.id`

	if activeOnly {
		query += ` WHERE d.is_active = true`
	}
	query += ` ORDER BY d.name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes []domain.Dish
	for rows.Next() {
		var dish domain.Dish
		dish.Category = &domain.Category{}

		if err := rows.Scan(
			&dish.ID, &dish.CategoryID, &dish.Name, &dish.Description,
			&dish.Price, &dish.PhotoURL, &dish.IsActive,
			&dish.Category.ID, &dish.Category.Name,
		); err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}

	return dishes, rows.Err()
}

func (r *DishRepository) GetByID(ctx context.Context, id int) (*domain.Dish, error) {
	query := `
		SELECT d.id, d.category_id, d.name, d.description, d.price, d.photo_url, d.is_active,
		       c.id, c.name
		FROM dishes d
		LEFT JOIN categories c ON d.category_id = c.id
		WHERE d.id = $1`

	dish := &domain.Dish{Category: &domain.Category{}}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dish.ID, &dish.CategoryID, &dish.Name, &dish.Description,
		&dish.Price, &dish.PhotoURL, &dish.IsActive,
		&dish.Category.ID, &dish.Category.Name,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return dish, err
}

func (r *DishRepository) GetByCategoryID(ctx context.Context, categoryID int) ([]domain.Dish, error) {
	query := `
		SELECT id, category_id, name, description, price, photo_url, is_active
		FROM dishes WHERE category_id = $1 AND is_active = true
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes []domain.Dish
	for rows.Next() {
		var dish domain.Dish
		if err := rows.Scan(
			&dish.ID, &dish.CategoryID, &dish.Name, &dish.Description,
			&dish.Price, &dish.PhotoURL, &dish.IsActive,
		); err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}

	return dishes, rows.Err()
}

func (r *DishRepository) Create(ctx context.Context, dish *domain.Dish) error {
	query := `
		INSERT INTO dishes (category_id, name, description, price, photo_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		dish.CategoryID, dish.Name, dish.Description, dish.Price, dish.PhotoURL, dish.IsActive,
	).Scan(&dish.ID)
}

func (r *DishRepository) Update(ctx context.Context, dish *domain.Dish) error {
	query := `
		UPDATE dishes 
		SET category_id = $1, name = $2, description = $3, price = $4, photo_url = $5, is_active = $6
		WHERE id = $7`

	_, err := r.db.ExecContext(ctx, query,
		dish.CategoryID, dish.Name, dish.Description, dish.Price, dish.PhotoURL, dish.IsActive, dish.ID,
	)
	return err
}

func (r *DishRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM dishes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *DishRepository) GetIngredients(ctx context.Context, dishID int) ([]domain.DishIngredient, error) {
	query := `
		SELECT di.dish_id, di.ingredient_id, di.qty_per_dish,
		       i.id, i.name, i.unit, i.qty, i.min_qty
		FROM dish_ingredients di
		JOIN ingredients i ON di.ingredient_id = i.id
		WHERE di.dish_id = $1`

	rows, err := r.db.QueryContext(ctx, query, dishID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []domain.DishIngredient
	for rows.Next() {
		var di domain.DishIngredient
		di.Ingredient = &domain.Ingredient{}

		if err := rows.Scan(
			&di.DishID, &di.IngredientID, &di.QtyPerDish,
			&di.Ingredient.ID, &di.Ingredient.Name, &di.Ingredient.Unit,
			&di.Ingredient.Qty, &di.Ingredient.MinQty,
		); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, di)
	}

	return ingredients, rows.Err()
}

func (r *DishRepository) AddIngredient(ctx context.Context, di *domain.DishIngredient) error {
	query := `
		INSERT INTO dish_ingredients (dish_id, ingredient_id, qty_per_dish)
		VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, di.DishID, di.IngredientID, di.QtyPerDish)
	return err
}

func (r *DishRepository) RemoveIngredient(ctx context.Context, dishID, ingredientID int) error {
	query := `DELETE FROM dish_ingredients WHERE dish_id = $1 AND ingredient_id = $2`
	_, err := r.db.ExecContext(ctx, query, dishID, ingredientID)
	return err
}

func (r *DishRepository) UpdateIngredient(ctx context.Context, di *domain.DishIngredient) error {
	query := `
		UPDATE dish_ingredients 
		SET qty_per_dish = $1
		WHERE dish_id = $2 AND ingredient_id = $3`

	_, err := r.db.ExecContext(ctx, query, di.QtyPerDish, di.DishID, di.IngredientID)
	return err
}
