// internal/adapters/postgres/category_repository.go
package postgre

import (
	"context"
	"database/sql"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, rows.Err()
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = $1`

	category := &domain.Category{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return category, err
}

func (r *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	return r.db.QueryRowContext(ctx, query, category.Name).Scan(&category.ID)
}

func (r *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `UPDATE categories SET name = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, category.Name, category.ID)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
