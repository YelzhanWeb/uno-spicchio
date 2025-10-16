// file: internal/controller/http/dto.go (или похожее место)

package httpAdapter

// CreateOrderItemRequest представляет одно блюдо в заказе
type CreateOrderItemRequest struct {
	DishID int    `json:"dish_id" validate:"required"`
	Qty    int    `json:"qty" validate:"required,gt=0"`
	Notes  string `json:"notes"`
}

// CreateOrderRequest - это структура для запроса на создание заказа
type CreateOrderRequest struct {
	WaiterID int                      `json:"waiter_id" validate:"required"`
	TableID  int                      `json:"table_id" validate:"required"`
	Notes    string                   `json:"notes"`
	Items    []CreateOrderItemRequest `json:"items" validate:"required,min=1"`
}
