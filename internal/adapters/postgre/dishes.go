// file: internal/adapters/postgre/dishes.go

package postgre

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/jmoiron/sqlx"
)

// DishRepository - это реализация порта DishRepository для PostgreSQL.
type DishRepository struct {
	db *sqlx.DB
}

// NewDishRepository создает новый экземпляр репозитория для блюд.
func NewDishRepository(db *sqlx.DB) *DishRepository {
	return &DishRepository{db: db}
}

// GetDishByID - получить блюдо по id.
func (r *DishRepository) GetDishByID(ctx context.Context, id int) (*domain.Dish, error) {
	var dish domain.Dish
	query := `SELECT id, category_id, name, description, price, photo_url, is_active
              FROM dishes WHERE id = $1`

	if err := r.db.GetContext(ctx, &dish, query, id); err != nil {
		return nil, fmt.Errorf("failed to get dish by id %d: %w", id, err)
	}
	return &dish, nil
}

// GetDishesByIDs - получить несколько блюд по списку их ID.
// Этот метод нужен для проверки цен при создании заказа.
func (r *DishRepository) GetDishesByIDs(ctx context.Context, ids []int) ([]domain.Dish, error) {
	if len(ids) == 0 {
		return []domain.Dish{}, nil
	}

	query, args, err := sqlx.In(`
		SELECT id, category_id, name, description, price, photo_url, is_active
		FROM dishes WHERE id IN (?)`, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to create IN query for dish ids: %w", err)
	}

	// sqlx.In возвращает query для конкретной СУБД, нужно его переделать под PostgreSQL.
	query = r.db.Rebind(query)

	var dishes []domain.Dish
	if err := r.db.SelectContext(ctx, &dishes, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get dishes by ids: %w", err)
	}

	return dishes, nil
}

// GetAllDishes - получить все блюда (оставил ваш код, адаптировав под sqlx)
func (r *DishRepository) GetAllDishes(ctx context.Context) ([]domain.Dish, error) {
	var dishes []domain.Dish
	query := `SELECT id, category_id, name, description, price, photo_url, is_active FROM dishes`
	if err := r.db.SelectContext(ctx, &dishes, query); err != nil {
		return nil, fmt.Errorf("failed to get all dishes: %w", err)
	}
	return dishes, nil
}

// Остальные ваши методы (Create, Update, Delete) можно добавить сюда по аналогии.

func (r *DishRepository) Create(ctx context.Context, dish domain.Dish) (int, error) {
	query := `INSERT INTO dishes (category_id, name, description, price, photo_url, is_active)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int
	err := r.db.QueryRowxContext(
		ctx,
		query,
		dish.CategoryID,
		dish.Name,
		dish.Description,
		dish.Price,
		dish.PhotoURL,
		dish.IsActive,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create dish: %w", err)
	}
	return id, nil
}
