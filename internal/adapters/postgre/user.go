package postgre

import (
	"context"
	"database/sql"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (username, password_hash, role, photokey, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query,
		user.Username, user.PasswordHash, user.Role, user.PhotoKey, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, username, password_hash, role, photokey, is_active, created_at
		FROM users WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Role,
		&user.PhotoKey, &user.IsActive, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, password_hash, role, photokey, is_active, created_at
		FROM users WHERE username = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Role,
		&user.PhotoKey, &user.IsActive, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	query := `
		SELECT id, username, password_hash, role, photokey, is_active, created_at
		FROM users ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID, &user.Username, &user.PasswordHash, &user.Role,
			&user.PhotoKey, &user.IsActive, &user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET username = $1, password_hash = $2, role = $3, photokey = $4, is_active = $5
		WHERE id = $6`

	_, err := r.db.ExecContext(ctx, query,
		user.Username, user.PasswordHash, user.Role, user.PhotoKey, user.IsActive, user.ID,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
