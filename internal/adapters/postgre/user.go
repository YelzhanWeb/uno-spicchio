package postgre

import "github.com/YelzhanWeb/uno-spicchio/internal/models"

func (r *Pool) CreateUser(username string, passwordHash string, role string) (int, error) {
	stmt := `INSERT INTO users (username, password_hash, role)
	         VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := r.DB.QueryRow(stmt, username, passwordHash, role).Scan(&id)
	return id, err
}

func (r *Pool) GetUserByID(id int) (*models.User, error) {
	stmt := `SELECT id, username, password_hash, role, created_at 
	         FROM users WHERE id = $1`
	row := r.DB.QueryRow(stmt, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Pool) GetAllUsers() ([]models.User, error) {
	stmt := `SELECT id, username, password_hash, role, created_at FROM users`
	rows, err := r.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *Pool) UpdateUser(id int, username string, passwordHash string, role string) error {
	stmt := `UPDATE users SET username=$1, password_hash=$2, role=$3 WHERE id=$4`
	_, err := r.DB.Exec(stmt, username, passwordHash, role, id)
	return err
}

func (r *Pool) DeleteUser(id int) error {
	stmt := `DELETE FROM users WHERE id=$1`
	_, err := r.DB.Exec(stmt, id)
	return err
}
