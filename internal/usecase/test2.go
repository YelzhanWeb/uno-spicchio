// // file: internal/usecase/usecase.go

package usecase

// import (
// 	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
// )

// // Dependencies - это структура для передачи всех репозиториев в слой usecase.
// type Dependencies struct {
// 	UserRepo     ports.UserRepository
// 	CategoryRepo ports.CategoryRepository
// 	OrderRepo    ports.OrderRepository
// 	DishRepo     ports.DishRepository
// 	// В будущем здесь будут и другие репозитории...
// }

// type UseCases struct {
// 	User     ports.UserService
// 	Category ports.CategoryService
// 	Order    ports.OrderService
// 	Dish     *DishService // <-- ДОБАВИТЬ ЭТУ СТРОКУ
// }

// // NewUseCases - это конструктор, который создает все наши сервисы.
// func NewUseCases(deps Dependencies) *UseCases {
// 	// Инициализируем каждый сервис, передавая ему нужный репозиторий
// 	userService := NewUserService(deps.UserRepo)
// 	categoryService := NewCategoryService(deps.CategoryRepo)
// 	orderService := NewOrderService(deps.OrderRepo, deps.DishRepo)
// 	dishService := NewDishService(deps.DishRepo)

// 	// Возвращаем контейнер со всеми сервисами
// 	return &UseCases{
// 		User:     userService,
// 		Category: categoryService,
// 		Order:    orderService,
// 		Dish:     dishService,
// 	}
// }
