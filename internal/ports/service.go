package ports

import (
	"context"
	"time"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

// AuthService defines methods for authentication
type AuthService interface {
	Login(ctx context.Context, username, password string) (string, *domain.User, error)
	GetCurrentUser(ctx context.Context, userID int) (*domain.User, error)
}

// UserService defines methods for user management
type UserService interface {
	Create(ctx context.Context, user *domain.User, password string) error
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, userID int, newPassword string) error
	Delete(ctx context.Context, id int) error
}

// OrderService defines methods for order management
type OrderService interface {
	Create(ctx context.Context, order *domain.Order, items []domain.OrderItem) error
	GetByID(ctx context.Context, id int) (*domain.Order, error)
	GetAll(ctx context.Context, status *domain.OrderStatus) ([]domain.Order, error)
	UpdateStatus(ctx context.Context, id int, newStatus domain.OrderStatus) error
	CloseOrder(ctx context.Context, id int) error
}

// DishService defines methods for dish management
type DishService interface {
	GetAll(ctx context.Context, activeOnly bool) ([]domain.Dish, error)
	GetByID(ctx context.Context, id int) (*domain.Dish, error)
	GetByCategoryID(ctx context.Context, categoryID int) ([]domain.Dish, error)
	Create(ctx context.Context, dish *domain.Dish) error
	Update(ctx context.Context, dish *domain.Dish) error
	Delete(ctx context.Context, id int) error
	GetIngredients(ctx context.Context, dishID int) ([]domain.DishIngredient, error)
	AddIngredient(ctx context.Context, dishIngredient *domain.DishIngredient) error
	RemoveIngredient(ctx context.Context, dishID, ingredientID int) error
}

// IngredientService defines methods for ingredient management
type IngredientService interface {
	GetAll(ctx context.Context) ([]domain.Ingredient, error)
	GetByID(ctx context.Context, id int) (*domain.Ingredient, error)
	GetLowStock(ctx context.Context) ([]domain.Ingredient, error)
	Create(ctx context.Context, ingredient *domain.Ingredient) error
	Update(ctx context.Context, ingredient *domain.Ingredient) error
	Delete(ctx context.Context, id int) error
}

// SupplyService defines methods for supply management
type SupplyService interface {
	Create(ctx context.Context, supply *domain.Supply) error
	GetAll(ctx context.Context) ([]domain.Supply, error)
	GetByIngredientID(ctx context.Context, ingredientID int) ([]domain.Supply, error)
}

// TableService defines methods for table management
type TableService interface {
	GetAll(ctx context.Context) ([]domain.Table, error)
	GetByID(ctx context.Context, id int) (*domain.Table, error)
	Create(ctx context.Context, table *domain.Table) error
	UpdateStatus(ctx context.Context, id int, status domain.TableStatus) error
	Delete(ctx context.Context, id int) error
}

type CategoryService interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetByID(ctx context.Context, id int) (*domain.Category, error)
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int) error
}

type AnalyticsService interface {
	GetDashboard(ctx context.Context, period domain.PeriodType, from, to time.Time) (*domain.DashboardData, error)
	GetSalesSummary(ctx context.Context, from, to time.Time) (*domain.SalesSummary, error)
	GetSalesByCategory(ctx context.Context, from, to time.Time) ([]domain.CategorySale, error)
	GetPopularDishes(ctx context.Context, from, to time.Time, limit int) ([]domain.PopularDish, error)
	GetWaiterPerformance(ctx context.Context, from, to time.Time) ([]domain.WaiterPerformance, error)
	GetOrderStats(ctx context.Context, from, to time.Time) (*domain.OrderStats, error)
	GetIngredientTurnover(ctx context.Context, from, to time.Time) ([]domain.IngredientTurnover, error)
	GetTableUtilization(ctx context.Context, from, to time.Time) ([]domain.TableUtilization, error)
	GetHourlyRevenue(ctx context.Context, date time.Time) ([]domain.HourlyRevenue, error)
}
