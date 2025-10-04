package postgre

import "github.com/YelzhanWeb/uno-spicchio/internal/domain"

func (r *Pool) CreateUser(username string, passwordHash string, role string, photoKey string) (int, error) {
	stmt := `INSERT INTO users (username, password_hash, role, photo_key)
	         VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.DB.QueryRow(stmt, username, passwordHash, role, photoKey).Scan(&id)
	return id, err
}

func (r *Pool) GetUserByID(id int) (*domain.User, error) {
	stmt := `SELECT id, username, password_hash, role, photo_key, created_at 
	         FROM users WHERE id = $1`
	row := r.DB.QueryRow(stmt, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.PhotoKey, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Pool) GetAllUsers() ([]domain.User, error) {
	stmt := `SELECT id, username, password_hash, role, photo_key, created_at FROM users`
	rows, err := r.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.PhotoKey, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *Pool) UpdateUser(id int, username string, passwordHash string, role string, photoKey string) error {
	stmt := `UPDATE users SET username=$1, password_hash=$2, role=$3, photo_key=$4 WHERE id=$5`
	_, err := r.DB.Exec(stmt, username, passwordHash, role, photoKey, id)
	return err
}

func (r *Pool) DeleteUser(id int) error {
	stmt := `DELETE FROM users WHERE id = $1`
	_, err := r.DB.Exec(stmt, id)
	return err
}
