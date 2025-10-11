// file: internal/ports/ports.go

package ports

import (
	"context"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
)

// =================================================================================
// РЕПОЗИТОРИИ (ИНТЕРФЕЙСЫ ДЛЯ РАБОТЫ С БАЗОЙ ДАННЫХ)
// =================================================================================

// UserRepository описывает методы для работы с пользователями в БД
type UserRepository interface {
	CreateUser(ctx context.Context, username, passwordHash, role, photoKey string) (int, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id int, username, passwordHash, role, photoKey string) error
	DeleteUser(ctx context.Context, id int) error
}

// CategoryRepository описывает методы для работы с категориями в БД
type CategoryRepository interface {
	CreateCategory(ctx context.Context, name string) (int, error)
	GetCategoryByID(ctx context.Context, id int) (*domain.Category, error)
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, id int, name string) error
	DeleteCategory(ctx context.Context, id int) error
	GetAllWithDishes(ctx context.Context) ([]domain.Category, error) // <-- МЕТОД ДЛЯ МЕНЮ ТЕПЕРЬ ЗДЕСЬ
}

// OrderRepository описывает методы для работы с заказами в БД
type OrderRepository interface {
	Create(ctx context.Context, order domain.Order) (int, error)
	UpdateStatus(ctx context.Context, orderID int, status string) error
	GetActiveWithItems(ctx context.Context) ([]domain.Order, error)
}

// DishRepository описывает методы для работы с блюдами в БД
type DishRepository interface {
	Create(ctx context.Context, dish domain.Dish) (int, error)
	GetDishByID(ctx context.Context, id int) (*domain.Dish, error)
	GetDishesByIDs(ctx context.Context, ids []int) ([]domain.Dish, error)
}

// =================================================================================
// СЕРВИСЫ (ИНТЕРФЕЙСЫ ДЛЯ БИЗНЕС-ЛОГИКИ / USE CASES)
// =================================================================================

// UserService описывает методы бизнес-логики для пользователей
type UserService interface {
	CreateUser(ctx context.Context, username, password, role, photoKey string) (int, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id int, username, password, role, photoKey string) error
	DeleteUser(ctx context.Context, id int) error
}

// CategoryService описывает методы бизнес-логики для категорий
type CategoryService interface {
	CreateCategory(ctx context.Context, name string) (int, error)
	GetCategoryByID(ctx context.Context, id int) (*domain.Category, error)
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, id int, name string) error
	DeleteCategory(ctx context.Context, id int) error
}

// OrderService описывает методы бизнес-логики для заказов
type OrderService interface {
	CreateOrder(ctx context.Context, order domain.Order) (int, error)
	UpdateOrderStatus(ctx context.Context, orderID int, status string) error
	GetActiveOrders(ctx context.Context) ([]domain.Order, error)
}

// HttpHandlers - оставляем без изменений
type HttpHandlers interface {
	PostUserHandler()
}
