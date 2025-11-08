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
	// Check stock for all items
	for _, item := range items {
		ingredients, err := s.dishRepo.GetIngredients(ctx, item.DishID)
		if err != nil {
			return fmt.Errorf("failed to get ingredients: %w", err)
		}

		for _, ing := range ingredients {
			needed := ing.QtyPerDish * float64(item.Qty)
			ingredient, err := s.ingredientRepo.GetByID(ctx, ing.IngredientID)
			if err != nil {
				return fmt.Errorf("failed to get ingredient: %w", err)
			}

			if ingredient.Qty < needed {
				return domain.ErrInsufficientStock
			}
		}
	}

	// Calculate total
	var total float64
	for _, item := range items {
		dish, err := s.dishRepo.GetByID(ctx, item.DishID)
		if err != nil {
			return fmt.Errorf("failed to get dish: %w", err)
		}
		total += dish.Price * float64(item.Qty)
	}

	order.Status = domain.OrderNew
	order.Total = total

	// Create order
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Add items
	for _, item := range items {
		item.OrderID = order.ID
		dish, _ := s.dishRepo.GetByID(ctx, item.DishID)
		item.Price = dish.Price

		if err := s.orderRepo.AddItem(ctx, &item); err != nil {
			return fmt.Errorf("failed to add order item: %w", err)
		}

		// Deduct ingredients
		ingredients, _ := s.dishRepo.GetIngredients(ctx, item.DishID)
		for _, ing := range ingredients {
			needed := ing.QtyPerDish * float64(item.Qty)
			if err := s.ingredientRepo.UpdateQuantity(ctx, ing.IngredientID, -needed); err != nil {
				return fmt.Errorf("failed to update ingredient quantity: %w", err)
			}
		}
	}

	// Update table status to BUSY
	if err := s.tableRepo.UpdateStatus(ctx, order.TableNumber, domain.TableBusy); err != nil {
		return fmt.Errorf("failed to update table status: %w", err)
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

	// Load items for each order
	for i := range orders {
		items, err := s.orderRepo.GetItems(ctx, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
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

	// Validate status transitions
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

	// Update order status to paid
	if err := s.orderRepo.UpdateStatus(ctx, id, domain.OrderPaid); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Update table status to FREE
	if err := s.tableRepo.UpdateStatus(ctx, order.TableNumber, domain.TableFree); err != nil {
		return fmt.Errorf("failed to update table status: %w", err)
	}

	return nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, id int) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}

	if err := s.orderRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}
