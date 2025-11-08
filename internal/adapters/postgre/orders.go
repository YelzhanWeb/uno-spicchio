package postgre

import (
	"context"
	"database/sql"
	"time"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (waiter_id, table_number, status, total, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		order.WaiterID, order.TableNumber, order.Status, order.Total, order.Notes,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (*domain.Order, error) {
	query := `
		SELECT 
			o.id, o.waiter_id, o.table_number, o.status, o.total, o.notes, 
			o.created_at, o.updated_at,
			u.id, u.username, u.role, u.photokey, u.is_active, u.created_at,
			t.id, t.name, t.status
		FROM orders o
		LEFT JOIN users u ON o.waiter_id = u.id
		LEFT JOIN tables t ON o.table_number = t.id
		WHERE o.id = $1`

	order := &domain.Order{
		Waiter: &domain.User{},
		Table:  &domain.Table{},
	}

	var waiterCreatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.WaiterID, &order.TableNumber, &order.Status, &order.Total, &order.Notes,
		&order.CreatedAt, &order.UpdatedAt,
		&order.Waiter.ID, &order.Waiter.Username, &order.Waiter.Role, &order.Waiter.PhotoKey,
		&order.Waiter.IsActive, &waiterCreatedAt,
		&order.Table.ID, &order.Table.Name, &order.Table.Status,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	order.Waiter.CreatedAt = waiterCreatedAt
	return order, nil
}

func (r *OrderRepository) GetAll(ctx context.Context, status *domain.OrderStatus) ([]domain.Order, error) {
	query := `
		SELECT 
			o.id, o.waiter_id, o.table_number, o.status, o.total, o.notes, 
			o.created_at, o.updated_at,
			COALESCE(u.username, '') as waiter_username,
			COALESCE(t.name, '') as table_name,
			t.id as table_id
		FROM orders o
		LEFT JOIN users u ON o.waiter_id = u.id
		LEFT JOIN tables t ON o.table_number = t.id`

	args := []interface{}{}
	if status != nil {
		query += ` WHERE o.status = $1`
		args = append(args, *status)
	}
	query += ` ORDER BY o.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var waiterUsername, tableName string
		var tableID int

		if err := rows.Scan(
			&order.ID, &order.WaiterID, &order.TableNumber, &order.Status, &order.Total,
			&order.Notes, &order.CreatedAt, &order.UpdatedAt,
			&waiterUsername, &tableName, &tableID,
		); err != nil {
			return nil, err
		}

		// Инициализируем связанные объекты
		order.Waiter = &domain.User{Username: waiterUsername}
		order.Table = &domain.Table{
			ID:   tableID,
			Name: tableName,
		}

		orders = append(orders, order)
	}

	return orders, rows.Err()
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	query := `
		UPDATE orders 
		SET waiter_id = $1, table_number = $2, status = $3, total = $4, notes = $5, updated_at = $6
		WHERE id = $7`

	_, err := r.db.ExecContext(ctx, query,
		order.WaiterID, order.TableNumber, order.Status, order.Total, order.Notes, time.Now(), order.ID,
	)
	return err
}

func (r *OrderRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *OrderRepository) AddItem(ctx context.Context, item *domain.OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, dish_id, qty, price, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		item.OrderID, item.DishID, item.Qty, item.Price, item.Notes,
	).Scan(&item.ID)
}

func (r *OrderRepository) GetItems(ctx context.Context, orderID int) ([]domain.OrderItem, error) {
	query := `
		SELECT 
			oi.id, oi.order_id, oi.dish_id, oi.qty, oi.price, 
			COALESCE(oi.notes, '') as notes,
			d.id, d.name, d.price, COALESCE(d.photo_url, '') as photo_url
		FROM order_items oi
		JOIN dishes d ON oi.dish_id = d.id
		WHERE oi.order_id = $1
		ORDER BY oi.id`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		var notes string
		item.Dish = &domain.Dish{}

		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.DishID, &item.Qty, &item.Price,
			&notes,
			&item.Dish.ID, &item.Dish.Name, &item.Dish.Price, &item.Dish.PhotoURL,
		); err != nil {
			return nil, err
		}

		if notes != "" {
			item.Notes = &notes
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *OrderRepository) UpdateItem(ctx context.Context, item *domain.OrderItem) error {
	query := `
		UPDATE order_items 
		SET dish_id = $1, qty = $2, price = $3, notes = $4
		WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query, item.DishID, item.Qty, item.Price, item.Notes, item.ID)
	return err
}

func (r *OrderRepository) DeleteItem(ctx context.Context, itemID int) error {
	query := `DELETE FROM order_items WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, itemID)
	return err
}
