// // Файл: internal/domain/orders.go
// package domain

// import "time"

// // Order представляет полную информацию о заказе.
// type Order struct {
// 	ID        int         `json:"id" db:"id"`
// 	WaiterID  int         `json:"waiter_id" db:"waiter_id"`
// 	TableID   int         `json:"table_id" db:"table_number"` // <-- Обратите внимание: колонка в БД 'table_number'
// 	Status    string      `json:"status" db:"status"`
// 	Total     float64     `json:"total" db:"total"`
// 	Notes     string      `json:"notes" db:"notes"`
// 	CreatedAt time.Time   `json:"created_at" db:"created_at"`
// 	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
// 	Items     []OrderItem `json:"items"` // Это поле не из БД, поэтому тег db не нужен
// }

// // OrderItem представляет одну позицию (блюдо) в заказе.
// // Убедитесь, что эта структура тоже имеет все нужные теги.
// type OrderItem struct {
// 	ID       int     `json:"id" db:"id"`
// 	OrderID  int     `json:"order_id" db:"order_id"`
// 	DishID   int     `json:"dish_id" db:"dish_id"`
// 	DishName string  `json:"dish_name" db:"dish_name"`
// 	Qty      int     `json:"qty" db:"qty"`
// 	Price    float64 `json:"price" db:"price"`
// 	Notes    string  `json:"notes" db:"notes"`
// }

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
