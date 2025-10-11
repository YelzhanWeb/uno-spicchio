// file: internal/usecase/orders.go

package usecase

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

// OrderService - это реализация бизнес-логики для заказов.
type OrderService struct {
	orderRepo ports.OrderRepository
	dishRepo  ports.DishRepository
}

// NewOrderService создает новый экземпляр сервиса заказов.
func NewOrderService(orderRepo ports.OrderRepository, dishRepo ports.DishRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		dishRepo:  dishRepo,
	}
}

// CreateOrder - основной метод бизнес-логики создания заказа.
func (s *OrderService) CreateOrder(ctx context.Context, order domain.Order) (int, error) {
	if len(order.Items) == 0 {
		return 0, fmt.Errorf("order must contain at least one item")
	}

	// 1. Собираем ID всех блюд из заказа, чтобы получить их из БД одним запросом.
	dishIDs := make([]int, len(order.Items))
	for i, item := range order.Items {
		dishIDs[i] = item.DishID
	}

	// 2. Получаем актуальные данные о блюдах (главное - цену) из БД.
	// Это ВАЖНО, чтобы клиент не мог отправить свою цену.
	dishes, err := s.dishRepo.GetDishesByIDs(ctx, dishIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to get dishes info: %w", err)
	}

	// Преобразуем слайс в мапу для быстрого доступа по ID.
	dishMap := make(map[int]domain.Dish)
	for _, dish := range dishes {
		dishMap[dish.ID] = dish
	}

	// 3. Проверяем, что все блюда найдены, и рассчитываем общую стоимость.
	var total float64
	for i := range order.Items {
		dish, ok := dishMap[order.Items[i].DishID]
		if !ok {
			return 0, fmt.Errorf("dish with id %d not found", order.Items[i].DishID)
		}
		// Устанавливаем цену из базы данных, а не от клиента.
		order.Items[i].Price = dish.Price
		total += dish.Price * float64(order.Items[i].Qty)
	}

	// 4. Устанавливаем рассчитанную сумму и статус для заказа.
	order.Total = total
	order.Status = "new" // Статус по умолчанию при создании.

	// 5. Вызываем репозиторий для сохранения заказа в БД.
	return s.orderRepo.Create(ctx, order)
}
func (s *OrderService) GetActiveOrders(ctx context.Context) ([]domain.Order, error) {
	// В будущем здесь может быть дополнительная бизнес-логика,
	// например, проверка прав доступа пользователя.
	return s.orderRepo.GetActiveWithItems(ctx)
}
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	// Здесь можно добавить логику валидации:
	// например, проверить, что переданный статус ('in_progress') является допустимым
	// или что нельзя перевести заказ из статуса 'paid' обратно в 'new'.
	// Пока просто вызываем репозиторий.
	return s.orderRepo.UpdateStatus(ctx, orderID, status)
}
