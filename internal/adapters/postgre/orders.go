// file: internal/adapters/postgre/orders.go

package postgre

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/jmoiron/sqlx"
)

// OrderRepository - это реализация порта OrderRepository для PostgreSQL.
type OrderRepository struct {
	db *sqlx.DB
}

// NewOrderRepository создает новый экземпляр репозитория для заказов.
func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create создает новую запись заказа и все его позиции в рамках одной транзакции.
func (r *OrderRepository) Create(ctx context.Context, order domain.Order) (int, error) {
	// Начинаем транзакцию
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Используем defer для отката транзакции в случае любой ошибки.
	// Если мы успешно выполним Commit в конце, Rollback не будет иметь эффекта.
	defer tx.Rollback()

	// 1. Создаем основную запись в таблице 'orders'
	// Обратите внимание, что мы не вставляем total, так как он должен быть рассчитан.
	// Статус по умолчанию будет 'new' согласно схеме БД.
	createOrderQuery := `
		INSERT INTO orders (waiter_id, table_number, notes, total)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var orderID int
	err = tx.QueryRowxContext(
		ctx,
		createOrderQuery,
		order.WaiterID,
		order.TableID,
		order.Notes,
		order.Total, // total должен быть рассчитан на слое usecase
	).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to create order record: %w", err)
	}

	// 2. Создаем записи для каждой позиции заказа в 'order_items'
	createItemQuery := `
		INSERT INTO order_items (order_id, dish_id, qty, price, notes)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, item := range order.Items {
		_, err := tx.ExecContext(
			ctx,
			createItemQuery,
			orderID, // Используем ID, полученный на предыдущем шаге
			item.DishID,
			item.Qty,
			item.Price, // Цена должна быть установлена на слое usecase
			item.Notes,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to create order item for dish %d: %w", item.DishID, err)
		}
	}

	// Если все прошло без ошибок, подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return orderID, nil
}
func (r *OrderRepository) GetActiveWithItems(ctx context.Context) ([]domain.Order, error) {
	// 1. Сначала получаем все активные заказы (все, кроме оплаченных)
	var orders []domain.Order
	queryOrders := `SELECT * FROM orders WHERE status != 'paid' ORDER BY created_at ASC`
	if err := r.db.SelectContext(ctx, &orders, queryOrders); err != nil {
		return nil, fmt.Errorf("failed to get active orders: %w", err)
	}

	// Если активных заказов нет, возвращаем пустой срез
	if len(orders) == 0 {
		return []domain.Order{}, nil
	}

	// 2. Собираем ID всех найденных заказов
	orderIDs := make([]int, len(orders))
	for i, order := range orders {
		orderIDs[i] = order.ID
	}

	// 3. Одним запросом получаем ВСЕ позиции для ВСЕХ найденных заказов,
	//    сразу объединяя с таблицей блюд, чтобы получить их названия.
	var items []domain.OrderItem
	queryItems, args, err := sqlx.In(`
		SELECT oi.*, d.name as dish_name 
		FROM order_items oi
		JOIN dishes d ON oi.dish_id = d.id
		WHERE oi.order_id IN (?)`, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create IN query for items: %w", err)
	}

	queryItems = r.db.Rebind(queryItems)
	if err := r.db.SelectContext(ctx, &items, queryItems, args...); err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	// 4. Распределяем найденные позиции по их заказам для удобного доступа
	itemsByOrderID := make(map[int][]domain.OrderItem)
	for _, item := range items {
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
	}

	// 5. Прикрепляем отсортированные позиции к соответствующим заказам
	for i := range orders {
		orderID := orders[i].ID
		if orderItems, ok := itemsByOrderID[orderID]; ok {
			orders[i].Items = orderItems
		}
	}

	return orders, nil
}
func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID int, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`

	res, err := r.db.ExecContext(ctx, query, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// Проверяем, что заказ с таким ID вообще существовал
	if rowsAffected == 0 {
		return fmt.Errorf("order with id %d not found", orderID)
	}

	return nil
}
