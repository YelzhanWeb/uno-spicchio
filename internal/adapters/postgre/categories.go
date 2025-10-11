// file: internal/adapters/postgre/categories.go

package postgre

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/jmoiron/sqlx"
)

// CategoryRepository - это реализация порта CategoryRepository для PostgreSQL.
type CategoryRepository struct {
	db *sqlx.DB
}

// NewCategoryRepository создает новый экземпляр репозитория для категорий.
func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// CreateCategory - добавление новой категории
func (r *CategoryRepository) CreateCategory(ctx context.Context, name string) (int, error) {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	var id int
	err := r.db.QueryRowxContext(ctx, query, name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create category: %w", err)
	}
	return id, nil
}

// GetCategoryByID - получить категорию по id
func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	var category domain.Category
	query := `SELECT id, name FROM categories WHERE id = $1`

	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id %d: %w", id, err)
	}
	return &category, nil
}

// GetAllCategories - получить все категории
func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	query := `SELECT id, name FROM categories`
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}
	return categories, nil
}

// UpdateCategory - обновить категорию по id
func (r *CategoryRepository) UpdateCategory(ctx context.Context, id int, name string) error {
	query := `UPDATE categories SET name=$1 WHERE id=$2`
	_, err := r.db.ExecContext(ctx, query, name, id)
	if err != nil {
		return fmt.Errorf("failed to update category with id %d: %w", id, err)
	}
	return nil
}

// DeleteCategory - удалить категорию по id
func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category with id %d: %w", id, err)
	}
	return nil
}
func (r *CategoryRepository) GetAllWithDishes(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	// 1. Получаем все категории
	if err := r.db.SelectContext(ctx, &categories, "SELECT * FROM categories ORDER BY name ASC"); err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	if len(categories) == 0 {
		return []domain.Category{}, nil
	}

	var dishes []domain.Dish
	// 2. Получаем все активные блюда
	if err := r.db.SelectContext(ctx, &dishes, "SELECT * FROM dishes WHERE is_active = true ORDER BY name ASC"); err != nil {
		return nil, fmt.Errorf("failed to get dishes: %w", err)
	}

	// 3. Распределяем блюда по категориям
	dishesByCategory := make(map[int][]domain.Dish)
	for _, dish := range dishes {
		dishesByCategory[dish.CategoryID] = append(dishesByCategory[dish.CategoryID], dish)
	}

	// 4. Прикрепляем блюда к их категориям
	for i := range categories {
		catID := categories[i].ID
		if categoryDishes, ok := dishesByCategory[catID]; ok {
			categories[i].Dishes = categoryDishes
		}
	}

	return categories, nil
}
