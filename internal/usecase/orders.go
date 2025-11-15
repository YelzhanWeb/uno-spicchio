package usecase

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

type OrderService struct {
	orderRepo      ports.OrderRepository
	dishRepo       ports.DishRepository
	ingredientRepo ports.IngredientRepository
	tableRepo      ports.TableRepository
}

func NewOrderService(
	orderRepo ports.OrderRepository,
	dishRepo ports.DishRepository,
	ingredientRepo ports.IngredientRepository,
	tableRepo ports.TableRepository,
) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		dishRepo:       dishRepo,
		ingredientRepo: ingredientRepo,
		tableRepo:      tableRepo,
	}
}

func (s *OrderService) Create(ctx context.Context, order *domain.Order, items []domain.OrderItem) error {
	// Проверяем существование стола
	table, err := s.tableRepo.GetByID(ctx, order.TableNumber)
	if err != nil {
		return err
	}
	if table == nil {
		return domain.ErrTableNotFound
	}

	// Проверяем наличие ингредиентов для всех блюд
	for _, item := range items {
		// Проверяем существование блюда
		dish, err := s.dishRepo.GetByID(ctx, item.DishID)
		if err != nil {
			return err
		}
		if dish == nil {
			return fmt.Errorf("dish with id %d not found", item.DishID)
		}

		// Получаем ингредиенты блюда
		ingredients, err := s.dishRepo.GetIngredients(ctx, item.DishID)
		if err != nil {
			return err
		}

		// Проверяем достаточность запасов
		for _, ing := range ingredients {
			needed := ing.QtyPerDish * float64(item.Qty)
			ingredient, err := s.ingredientRepo.GetByID(ctx, ing.IngredientID)
			if err != nil {
				return err
			}

			if ingredient.Qty < needed {
				return domain.ErrInsufficientStock
			}
		}
	}

	// Подсчитываем общую сумму заказа
	var total float64
	for _, item := range items {
		dish, err := s.dishRepo.GetByID(ctx, item.DishID)
		if err != nil {
			return err
		}
		total += dish.Price * float64(item.Qty)
	}

	order.Status = domain.OrderNew
	order.Total = total

	// Создаем заказ
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return err
	}

	// Добавляем позиции заказа
	for _, item := range items {
		item.OrderID = order.ID
		dish, _ := s.dishRepo.GetByID(ctx, item.DishID)
		item.Price = dish.Price

		if err := s.orderRepo.AddItem(ctx, &item); err != nil {
			return err
		}

		// Списываем ингредиенты
		ingredients, _ := s.dishRepo.GetIngredients(ctx, item.DishID)
		for _, ing := range ingredients {
			needed := ing.QtyPerDish * float64(item.Qty)
			if err := s.ingredientRepo.UpdateQuantity(ctx, ing.IngredientID, -needed); err != nil {
				return err
			}
		}
	}

	// Обновляем статус стола на "занят"
	if err := s.tableRepo.UpdateStatus(ctx, order.TableNumber, domain.TableBusy); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetByID(ctx context.Context, id int) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}

	// Загружаем позиции заказа
	items, err := s.orderRepo.GetItems(ctx, id)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (s *OrderService) GetAll(ctx context.Context, status *domain.OrderStatus) ([]domain.Order, error) {
	orders, err := s.orderRepo.GetAll(ctx, status)
	if err != nil {
		return nil, err
	}

	// Загружаем детали для каждого заказа
	for i := range orders {
		items, err := s.orderRepo.GetItems(ctx, orders[i].ID)
		if err == nil {
			orders[i].Items = items
		}
	}

	return orders, nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, id int, newStatus domain.OrderStatus) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}

	// Валидация переходов статусов
	validTransitions := map[domain.OrderStatus][]domain.OrderStatus{
		domain.OrderNew:        {domain.OrderInProgress},
		domain.OrderInProgress: {domain.OrderReady},
		domain.OrderReady:      {domain.OrderPaid},
	}

	valid := false
	for _, allowedStatus := range validTransitions[order.Status] {
		if allowedStatus == newStatus {
			valid = true
			break
		}
	}

	// Разрешаем не менять статус, если он уже установлен
	if !valid && newStatus != order.Status {
		return domain.ErrInvalidStatusChange
	}

	return s.orderRepo.UpdateStatus(ctx, id, newStatus)
}

func (s *OrderService) CloseOrder(ctx context.Context, id int) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}

	// Проверяем, что заказ в статусе ready
	if order.Status != domain.OrderReady {
		return fmt.Errorf("order must be in ready status to close, current status: %s", order.Status)
	}

	// Обновляем статус заказа на "оплачен"
	if err := s.orderRepo.UpdateStatus(ctx, id, domain.OrderPaid); err != nil {
		return err
	}

	// Освобождаем стол
	if err := s.tableRepo.UpdateStatus(ctx, order.TableNumber, domain.TableFree); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) Delete(ctx context.Context, id int) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}

	// Можно добавить проверку: разрешить удалять только новые заказы
	if order.Status != domain.OrderNew {
		return fmt.Errorf("cannot delete order in status: %s", order.Status)
	}

	return s.orderRepo.Delete(ctx, id)
}
