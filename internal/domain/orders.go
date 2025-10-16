package domain

import "time"

type OrderStatus string

const (
	OrderNew        OrderStatus = "new"
	OrderInProgress OrderStatus = "in_progress"
	OrderReady      OrderStatus = "ready"
	OrderPaid       OrderStatus = "paid"
)

type Order struct {
	ID          int         `json:"id"`
	WaiterID    int         `json:"waiter_id"`
	TableNumber int         `json:"table_number"`
	Status      OrderStatus `json:"status"`
	Total       float64     `json:"total"`
	Notes       *string     `json:"notes,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`

	// Relations
	Items  []OrderItem `json:"items,omitempty"`
	Waiter *User       `json:"waiter,omitempty"`
	Table  *Table      `json:"table,omitempty"`
}

type OrderItem struct {
	ID      int     `json:"id"`
	OrderID int     `json:"order_id"`
	DishID  int     `json:"dish_id"`
	Qty     int     `json:"qty"`
	Price   float64 `json:"price"`
	Notes   *string `json:"notes,omitempty"`
	Dish    *Dish   `json:"dish,omitempty"`
}
