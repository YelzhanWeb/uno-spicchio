package usecase

import (
	"context"
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/pkg/logger"
)

type OrderService struct {
	orderRepo      ports.OrderRepository
	dishRepo       ports.DishRepository
	ingredientRepo ports.IngredientRepository
	tableRepo      ports.TableRepository
	logger         *logger.Logger
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
		logger:         logger.New("OrderService"),
	}
}

func (s *OrderService) Create(ctx context.Context, order *domain.Order, items []domain.OrderItem) error {
	s.logger.Order("Creating new order for table #%d", order.TableNumber)

	// Проверяем существование стола
	table, err := s.tableRepo.GetByID(ctx, order.TableNumber)
	if err != nil {
		s.logger.Error("Failed to get table #%d: %v", order.TableNumber, err)
		return err
	}
	if table == nil {
		s.logger.Error("Table #%d not found", order.TableNumber)
		return domain.ErrTableNotFound
	}

	s.logger.Info("Table #%d found: %s", table.ID, table.Name)

	// Проверяем наличие ингредиентов для всех блюд
	for _, item := range items {
		// Проверяем существование блюда
		dish, err := s.dishRepo.GetByID(ctx, item.DishID)
		if err != nil {
			s.logger.Error("Failed to get dish #%d: %v", item.DishID, err)
			return err
		}
		if dish == nil {
			s.logger.Error("Dish #%d not found", item.DishID)
			return fmt.Errorf("dish with id %d not found", item.DishID)
		}

		s.logger.Info("Adding dish '%s' (x%d) to order", dish.Name, item.Qty)

		// Получаем ингредиенты блюда
		ingredients, err := s.dishRepo.GetIngredients(ctx, item.DishID)
		if err != nil {
			s.logger.Error("Failed to get ingredients for dish #%d: %v", item.DishID, err)
			return err
		}

		// Проверяем достаточность запасов
		for _, ing := range ingredients {
			needed := ing.QtyPerDish * float64(item.Qty)
			ingredient, err := s.ingredientRepo.GetByID(ctx, ing.IngredientID)
			if err != nil {
				s.logger.Error("Failed to get ingredient #%d: %v", ing.IngredientID, err)
				return err
			}

			if ingredient.Qty < needed {
				s.logger.Error("Insufficient stock for ingredient '%s': needed %.2f, available %.2f",
					ingredient.Name, needed, ingredient.Qty)
				return domain.ErrInsufficientStock
			}

			s.logger.Debug("Ingredient '%s': needed %.2f%s, available %.2f%s",
				ingredient.Name, needed, ingredient.Unit, ingredient.Qty, ingredient.Unit)
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

	s.logger.Order("Order total: %.2f ₸", total)

	// Создаем заказ
	if err := s.orderRepo.Create(ctx, order); err != nil {
		s.logger.Error("Failed to create order: %v", err)
		return err
	}

	s.logger.Success("✓ Order #%d created successfully", order.ID)

	// Добавляем позиции заказа
	for _, item := range items {
		item.OrderID = order.ID
		dish, _ := s.dishRepo.GetByID(ctx, item.DishID)
		item.Price = dish.Price

		if err := s.orderRepo.AddItem(ctx, &item); err != nil {
			s.logger.Error("Failed to add item to order: %v", err)
			return err
		}

		s.logger.Info("Added item: %s (x%d) - %.2f ₸", dish.Name, item.Qty, item.Price)

		// Списываем ингредиенты
		ingredients, _ := s.dishRepo.GetIngredients(ctx, item.DishID)
		for _, ing := range ingredients {
			needed := ing.QtyPerDish * float64(item.Qty)
			if err := s.ingredientRepo.UpdateQuantity(ctx, ing.IngredientID, -needed); err != nil {
				s.logger.Error("Failed to update ingredient quantity: %v", err)
				return err
			}
			s.logger.Debug("Deducted %.2f from ingredient #%d", needed, ing.IngredientID)
		}
	}

	// Обновляем статус стола на "занят"
	if err := s.tableRepo.UpdateStatus(ctx, order.TableNumber, domain.TableBusy); err != nil {
		s.logger.Error("Failed to update table status: %v", err)
		return err
	}

	s.logger.Success("✓ Table #%d marked as busy", order.TableNumber)
	s.logger.Order("Order #%d created successfully with %d items", order.ID, len(items))

	return nil
}

func (s *OrderService) GetByID(ctx context.Context, id int) (*domain.Order, error) {
	s.logger.Info("Fetching order #%d", id)

	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get order #%d: %v", id, err)
		return nil, err
	}
	if order == nil {
		s.logger.Warning("Order #%d not found", id)
		return nil, domain.ErrOrderNotFound
	}

	// Загружаем позиции заказа
	items, err := s.orderRepo.GetItems(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get items for order #%d: %v", id, err)
		return nil, err
	}
	order.Items = items

	s.logger.Success("✓ Order #%d fetched successfully (%d items)", id, len(items))
	return order, nil
}

func (s *OrderService) GetAll(ctx context.Context, status *domain.OrderStatus) ([]domain.Order, error) {
	if status != nil {
		s.logger.Info("Fetching all orders with status: %s", *status)
	} else {
		s.logger.Info("Fetching all orders")
	}

	orders, err := s.orderRepo.GetAll(ctx, status)
	if err != nil {
		s.logger.Error("Failed to get orders: %v", err)
		return nil, err
	}

	// Загружаем детали для каждого заказа
	for i := range orders {
		items, err := s.orderRepo.GetItems(ctx, orders[i].ID)
		if err == nil {
			orders[i].Items = items
		}
	}

	s.logger.Success("✓ Fetched %d orders", len(orders))
	return orders, nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, id int, newStatus domain.OrderStatus) error {
	s.logger.Order("Updating order #%d status to: %s", id, newStatus)

	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get order #%d: %v", id, err)
		return err
	}
	if order == nil {
		s.logger.Error("Order #%d not found", id)
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
		s.logger.Error("Invalid status transition from %s to %s", order.Status, newStatus)
		return domain.ErrInvalidStatusChange
	}

	if err := s.orderRepo.UpdateStatus(ctx, id, newStatus); err != nil {
		s.logger.Error("Failed to update status: %v", err)
		return err
	}

	s.logger.Success("✓ Order #%d status updated: %s → %s", id, order.Status, newStatus)
	return nil
}

func (s *OrderService) CloseOrder(ctx context.Context, id int) error {
	s.logger.Order("Closing order #%d", id)

	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get order #%d: %v", id, err)
		return err
	}
	if order == nil {
		s.logger.Error("Order #%d not found", id)
		return domain.ErrOrderNotFound
	}

	// Проверяем, что заказ в статусе ready
	if order.Status != domain.OrderReady {
		s.logger.Error("Order #%d must be in 'ready' status to close, current: %s", id, order.Status)
		return fmt.Errorf("order must be in ready status to close, current status: %s", order.Status)
	}

	// Обновляем статус заказа на "оплачен"
	if err := s.orderRepo.UpdateStatus(ctx, id, domain.OrderPaid); err != nil {
		s.logger.Error("Failed to update order status: %v", err)
		return err
	}

	s.logger.Success("✓ Order #%d marked as paid", id)

	// Освобождаем стол
	if err := s.tableRepo.UpdateStatus(ctx, order.TableNumber, domain.TableFree); err != nil {
		s.logger.Error("Failed to free table #%d: %v", order.TableNumber, err)
		return err
	}

	s.logger.Success("✓ Table #%d freed", order.TableNumber)
	s.logger.Order("Order #%d closed successfully (Total: %.2f ₸)", id, order.Total)

	return nil
}

func (s *OrderService) Delete(ctx context.Context, id int) error {
	s.logger.Warning("Deleting order #%d", id)

	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get order #%d: %v", id, err)
		return err
	}
	if order == nil {
		s.logger.Error("Order #%d not found", id)
		return domain.ErrOrderNotFound
	}

	// Можно добавить проверку: разрешить удалять только новые заказы
	if order.Status != domain.OrderNew {
		s.logger.Error("Cannot delete order #%d in status: %s", id, order.Status)
		return fmt.Errorf("cannot delete order in status: %s", order.Status)
	}

	if err := s.orderRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete order #%d: %v", id, err)
		return err
	}

	s.logger.Success("✓ Order #%d deleted", id)
	return nil
}
