// file: internal/adapters/postgre/user.go

package postgre

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/jmoiron/sqlx"
)

// UserRepository - это реализация репозитория для пользователей.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository создает новый экземпляр репозитория для пользователей.
// Этот конструктор теперь будет найден в main.go!
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, username string, passwordHash string, role string, photoKey string) (int, error) {
	query := `INSERT INTO users (username, password_hash, role, photo_key)
	          VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.db.QueryRowxContext(ctx, query, username, passwordHash, role, photoKey).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	// Убедитесь, что в domain.User поля помечены тегами `db:"..."`, например Password string `db:"password_hash"`
	query := `SELECT id, username, password_hash, role, photo_key, is_active, created_at
	          FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id %d: %w", id, err)
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	query := `SELECT id, username, password_hash, role, photo_key, is_active, created_at FROM users`
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int, username string, passwordHash string, role string, photoKey string) error {
	query := `UPDATE users SET username=$1, password_hash=$2, role=$3, photo_key=$4 WHERE id=$5`
	_, err := r.db.ExecContext(ctx, query, username, passwordHash, role, photoKey, id)
	if err != nil {
		return fmt.Errorf("failed to update user with id %d: %w", id, err)
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user with id %d: %w", id, err)
	}
	return nil
}
