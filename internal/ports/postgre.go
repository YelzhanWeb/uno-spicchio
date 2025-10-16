package ports

import (
	"context"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

// UserRepository defines methods for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
}

// TableRepository defines methods for table data access
type TableRepository interface {
	GetAll(ctx context.Context) ([]domain.Table, error)
	GetByID(ctx context.Context, id int) (*domain.Table, error)
	Create(ctx context.Context, table *domain.Table) error
	UpdateStatus(ctx context.Context, id int, status domain.TableStatus) error
	Delete(ctx context.Context, id int) error
}

// CategoryRepository defines methods for category data access
type CategoryRepository interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetByID(ctx context.Context, id int) (*domain.Category, error)
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int) error
}

// DishRepository defines methods for dish data access
type DishRepository interface {
	GetAll(ctx context.Context, activeOnly bool) ([]domain.Dish, error)
	GetByID(ctx context.Context, id int) (*domain.Dish, error)
	GetByCategoryID(ctx context.Context, categoryID int) ([]domain.Dish, error)
	Create(ctx context.Context, dish *domain.Dish) error
	Update(ctx context.Context, dish *domain.Dish) error
	Delete(ctx context.Context, id int) error
	GetIngredients(ctx context.Context, dishID int) ([]domain.DishIngredient, error)
	AddIngredient(ctx context.Context, dishIngredient *domain.DishIngredient) error
	RemoveIngredient(ctx context.Context, dishID, ingredientID int) error
	UpdateIngredient(ctx context.Context, dishIngredient *domain.DishIngredient) error
}

// IngredientRepository defines methods for ingredient data access
type IngredientRepository interface {
	GetAll(ctx context.Context) ([]domain.Ingredient, error)
	GetByID(ctx context.Context, id int) (*domain.Ingredient, error)
	GetLowStock(ctx context.Context) ([]domain.Ingredient, error)
	Create(ctx context.Context, ingredient *domain.Ingredient) error
	Update(ctx context.Context, ingredient *domain.Ingredient) error
	UpdateQuantity(ctx context.Context, id int, qty float64) error
	Delete(ctx context.Context, id int) error
}

// OrderRepository defines methods for order data access
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id int) (*domain.Order, error)
	GetAll(ctx context.Context, status *domain.OrderStatus) ([]domain.Order, error)
	UpdateStatus(ctx context.Context, id int, status domain.OrderStatus) error
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, id int) error

	// Order items
	AddItem(ctx context.Context, item *domain.OrderItem) error
	GetItems(ctx context.Context, orderID int) ([]domain.OrderItem, error)
	UpdateItem(ctx context.Context, item *domain.OrderItem) error
	DeleteItem(ctx context.Context, itemID int) error
}

// SupplyRepository defines methods for supply data access
type SupplyRepository interface {
	Create(ctx context.Context, supply *domain.Supply) error
	GetAll(ctx context.Context) ([]domain.Supply, error)
	GetByID(ctx context.Context, id int) (*domain.Supply, error)
	GetByIngredientID(ctx context.Context, ingredientID int) ([]domain.Supply, error)
}
