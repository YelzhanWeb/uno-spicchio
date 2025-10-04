package postgre

import (
	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

// CreateOrder - создать новый заказ
func (r *Pool) CreateOrder(waiterID, tableNumber int, status string) (int, error) {
	stmt := `INSERT INTO orders (waiter_id, table_number, status, total)
             VALUES ($1, $2, $3, 0) RETURNING id`
	var id int
	err := r.DB.QueryRow(stmt, waiterID, tableNumber, status).Scan(&id)
	return id, err
}

// GetOrderByID - получить заказ по id (без позиций)
func (r *Pool) GetOrderByID(id int) (*domain.Order, error) {
	stmt := `SELECT id, waiter_id, table_number, status, total, created_at 
             FROM orders WHERE id=$1`
	row := r.DB.QueryRow(stmt, id)

	var o domain.Order
	err := row.Scan(&o.ID, &o.WaiterID, &o.TableNumber, &o.Status, &o.Total, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetAllOrders - получить все заказы (без позиций)
func (r *Pool) GetAllOrders() ([]domain.Order, error) {
	stmt := `SELECT id, waiter_id, table_number, status, total, created_at FROM orders`
	rows, err := r.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.WaiterID, &o.TableNumber, &o.Status, &o.Total, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// UpdateOrderStatus - обновить статус заказа
func (r *Pool) UpdateOrderStatus(id int, status string) error {
	stmt := `UPDATE orders SET status=$1 WHERE id=$2`
	_, err := r.DB.Exec(stmt, status, id)
	return err
}

// DeleteOrder - удалить заказ (позиции удалятся каскадно)
func (r *Pool) DeleteOrder(id int) error {
	stmt := `DELETE FROM orders WHERE id=$1`
	_, err := r.DB.Exec(stmt, id)
	return err
}

//////////////////////////////////////////////////
// Работа с order_items
//////////////////////////////////////////////////

// AddItem - добавить позицию в заказ
func (r *Pool) AddItem(orderID, dishID, qty int, price float64) (int, error) {
	stmt := `INSERT INTO order_items (order_id, dish_id, qty, price)
             VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.DB.QueryRow(stmt, orderID, dishID, qty, price).Scan(&id)
	if err != nil {
		return 0, err
	}

	// обновляем сумму заказа
	_, _ = r.DB.Exec(`UPDATE orders SET total = (
        SELECT COALESCE(SUM(qty * price),0) FROM order_items WHERE order_id=$1
    ) WHERE id=$1`, orderID)

	return id, nil
}

// UpdateItemQty - изменить количество позиции
func (r *Pool) UpdateItemQty(itemID, qty int) error {
	stmt := `UPDATE order_items SET qty=$1 WHERE id=$2`
	_, err := r.DB.Exec(stmt, qty, itemID)
	if err != nil {
		return err
	}

	// пересчитать total
	_, _ = r.DB.Exec(`UPDATE orders SET total = (
        SELECT COALESCE(SUM(qty * price),0) FROM order_items WHERE order_id=(
            SELECT order_id FROM order_items WHERE id=$1
        )
    ) WHERE id=(SELECT order_id FROM order_items WHERE id=$1)`, itemID)

	return nil
}

// RemoveItem - удалить позицию из заказа
func (r *Pool) RemoveItem(itemID int) error {
	var orderID int
	err := r.DB.QueryRow(`SELECT order_id FROM order_items WHERE id=$1`, itemID).Scan(&orderID)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(`DELETE FROM order_items WHERE id=$1`, itemID)
	if err != nil {
		return err
	}

	// пересчитать total
	_, _ = r.DB.Exec(`UPDATE orders SET total = (
        SELECT COALESCE(SUM(qty * price),0) FROM order_items WHERE order_id=$1
    ) WHERE id=$1`, orderID)

	return nil
}

// GetOrderWithItems - получить заказ вместе с позициями
func (r *Pool) GetOrderWithItems(orderID int) (*domain.OrderWithItems, error) {
	// сам заказ
	order, err := r.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}

	// его позиции
	rows, err := r.DB.Query(`SELECT id, order_id, dish_id, qty, price FROM order_items WHERE order_id=$1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var oi domain.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.DishID, &oi.Qty, &oi.Price); err != nil {
			return nil, err
		}
		items = append(items, oi)
	}

	return &domain.OrderWithItems{
		Order: *order,
		Items: items,
	}, nil
}
