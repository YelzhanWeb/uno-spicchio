package postgre

import (
	"context"
	"database/sql"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

type TableRepository struct {
	db *sql.DB
}

func NewTableRepository(db *sql.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) GetAll(ctx context.Context) ([]domain.Table, error) {
	query := `SELECT id, name, status FROM tables ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []domain.Table
	for rows.Next() {
		var table domain.Table
		if err := rows.Scan(&table.ID, &table.Name, &table.Status); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, rows.Err()
}

func (r *TableRepository) GetByID(ctx context.Context, id int) (*domain.Table, error) {
	query := `SELECT id, name, status FROM tables WHERE id = $1`

	table := &domain.Table{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&table.ID, &table.Name, &table.Status)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return table, err
}

func (r *TableRepository) Create(ctx context.Context, table *domain.Table) error {
	query := `INSERT INTO tables (name, status) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, table.Name, table.Status).Scan(&table.ID)
}

func (r *TableRepository) UpdateStatus(ctx context.Context, id int, status domain.TableStatus) error {
	query := `UPDATE tables SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *TableRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tables WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
