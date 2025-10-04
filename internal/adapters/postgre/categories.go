package postgre

import "github.com/YelzhanWeb/uno-spicchio/internal/domain"

// CreateCategory - добавление новой категории
func (r *Pool) CreateCategory(name string) (int, error) {
	stmt := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	var id int
	err := r.DB.QueryRow(stmt, name).Scan(&id)
	return id, err
}

// GetCategoryByID - получить категорию по id
func (r *Pool) GetCategoryByID(id int) (*domain.Category, error) {
	stmt := `SELECT id, name FROM categories WHERE id = $1`
	row := r.DB.QueryRow(stmt, id)

	var category domain.Category
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAllCategories - получить все категории
func (r *Pool) GetAllCategories() ([]domain.Category, error) {
	stmt := `SELECT id, name FROM categories`
	rows, err := r.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// UpdateCategory - обновить категорию по id
func (r *Pool) UpdateCategory(id int, name string) error {
	stmt := `UPDATE categories SET name=$1 WHERE id=$2`
	_, err := r.DB.Exec(stmt, name, id)
	return err
}

// DeleteCategory - удалить категорию по id
func (r *Pool) DeleteCategory(id int) error {
	stmt := `DELETE FROM categories WHERE id=$1`
	_, err := r.DB.Exec(stmt, id)
	return err
}
