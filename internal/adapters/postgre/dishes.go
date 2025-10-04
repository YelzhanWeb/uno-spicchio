package postgre

import "github.com/YelzhanWeb/uno-spicchio/internal/domain"

// CreateDish - добавить новое блюдо
func (r *Pool) CreateDish(categoryID int, name, description string, price float64, photoURL string) (int, error) {
	stmt := `INSERT INTO dishes (category_id, name, description, price, photo_url)
	         VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := r.DB.QueryRow(stmt, categoryID, name, description, price, photoURL).Scan(&id)
	return id, err
}

// GetDishByID - получить блюдо по id
func (r *Pool) GetDishByID(id int) (*domain.Dish, error) {
	stmt := `SELECT id, category_id, name, description, price, photo_url 
	         FROM dishes WHERE id = $1`
	row := r.DB.QueryRow(stmt, id)

	var dish domain.Dish
	err := row.Scan(&dish.ID, &dish.CategoryID, &dish.Name, &dish.Description, &dish.Price, &dish.PhotoURL)
	if err != nil {
		return nil, err
	}
	return &dish, nil
}

// GetAllDishes - получить все блюда
func (r *Pool) GetAllDishes() ([]domain.Dish, error) {
	stmt := `SELECT id, category_id, name, description, price, photo_url FROM dishes`
	rows, err := r.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes []domain.Dish
	for rows.Next() {
		var d domain.Dish
		if err := rows.Scan(&d.ID, &d.CategoryID, &d.Name, &d.Description, &d.Price, &d.PhotoURL); err != nil {
			return nil, err
		}
		dishes = append(dishes, d)
	}
	return dishes, nil
}

// UpdateDish - обновить блюдо по id
func (r *Pool) UpdateDish(id int, categoryID int, name, description string, price float64, photoURL string) error {
	stmt := `UPDATE dishes 
	         SET category_id=$1, name=$2, description=$3, price=$4, photo_url=$5 
	         WHERE id=$6`
	_, err := r.DB.Exec(stmt, categoryID, name, description, price, photoURL, id)
	return err
}

// DeleteDish - удалить блюдо по id
func (r *Pool) DeleteDish(id int) error {
	stmt := `DELETE FROM dishes WHERE id=$1`
	_, err := r.DB.Exec(stmt, id)
	return err
}
