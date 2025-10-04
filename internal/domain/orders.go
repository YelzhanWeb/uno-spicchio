package domain

import "time"

type Order struct {
	ID          int       `json:"id"`
	WaiterID    int       `json:"waiter_id"`
	TableNumber int       `json:"table_number"`
	Status      string    `json:"status"`
	Total       float64   `json:"total"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderItem struct {
	ID      int     `json:"id"`
	OrderID int     `json:"order_id"`
	DishID  int     `json:"dish_id"`
	Qty     int     `json:"qty"`
	Price   float64 `json:"price"`
}

// Заказ вместе с позициями
type OrderWithItems struct {
	Order Order       `json:"order"`
	Items []OrderItem `json:"items"`
}
