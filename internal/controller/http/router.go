// internal/adapters/http/router.go
package httpAdapter

import (
	"github.com/YelzhanWeb/uno-spicchio/internal/controller/http/handlers"
	"github.com/YelzhanWeb/uno-spicchio/internal/controller/http/middleware"
	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	authHandler       *handlers.AuthHandler
	userHandler       *handlers.UserHandler
	orderHandler      *handlers.OrderHandler
	dishHandler       *handlers.DishHandler
	ingredientHandler *handlers.IngredientHandler
	supplyHandler     *handlers.SupplyHandler
	tableHandler      *handlers.TableHandler
	categoryHandler   *handlers.CategoryHandler
	analyticsHandler  *handlers.AnalyticsHandler
	tokenManager      *jwt.TokenManager
}

func NewRouter(
	authService ports.AuthService,
	userService ports.UserService,
	orderService ports.OrderService,
	dishService ports.DishService,
	ingredientService ports.IngredientService,
	supplyService ports.SupplyService,
	tableService ports.TableService,
	categoryService ports.CategoryService,
	analyticsService ports.AnalyticsService,
	tokenManager *jwt.TokenManager,
) *Router {
	return &Router{
		authHandler:       handlers.NewAuthHandler(authService),
		userHandler:       handlers.NewUserHandler(userService),
		orderHandler:      handlers.NewOrderHandler(orderService),
		dishHandler:       handlers.NewDishHandler(dishService),
		ingredientHandler: handlers.NewIngredientHandler(ingredientService),
		supplyHandler:     handlers.NewSupplyHandler(supplyService),
		tableHandler:      handlers.NewTableHandler(tableService),
		categoryHandler:   handlers.NewCategoryHandler(categoryService),
		analyticsHandler:  handlers.NewAnalyticsHandler(analyticsService),
		tokenManager:      tokenManager,
	}
}

func (rt *Router) Setup() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logging)
	r.Use(middleware.CORS)

	// Public routes
	r.Post("/api/auth/login", rt.authHandler.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(rt.tokenManager))

		// Auth routes
		r.Get("/api/auth/me", rt.authHandler.GetMe)

		// User routes (Admin only)
		r.Route("/api/users", func(r chi.Router) {
			r.Use(middleware.RequireRole(domain.RoleAdmin))
			r.Get("/", rt.userHandler.GetAll)
			r.Post("/", rt.userHandler.Create)
			r.Get("/{id}", rt.userHandler.GetByID)
			r.Put("/{id}", rt.userHandler.Update)
			r.Delete("/{id}", rt.userHandler.Delete)
		})

		// Category routes
		r.Route("/api/categories", func(r chi.Router) {
			r.Get("/", rt.categoryHandler.GetAll)
			r.Get("/{id}", rt.categoryHandler.GetByID)

			// Admin only
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole(domain.RoleAdmin))
				r.Post("/", rt.categoryHandler.Create)
				r.Put("/{id}", rt.categoryHandler.Update)
				r.Delete("/{id}", rt.categoryHandler.Delete)
			})
		})

		// Dish routes
		r.Route("/api/dishes", func(r chi.Router) {
			r.Get("/", rt.dishHandler.GetAll)
			r.Get("/{id}", rt.dishHandler.GetByID)
			r.Get("/{id}/ingredients", rt.dishHandler.GetIngredients)

			// Admin only
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole(domain.RoleAdmin))
				r.Post("/", rt.dishHandler.Create)
				r.Put("/{id}", rt.dishHandler.Update)
				r.Delete("/{id}", rt.dishHandler.Delete)
			})
		})

		// Order routes
		r.Route("/api/orders", func(r chi.Router) {
			r.Get("/", rt.orderHandler.GetAll)
			r.Get("/{id}", rt.orderHandler.GetByID)

			// Waiter can create orders
			r.With(middleware.RequireRole(domain.RoleWaiter, domain.RoleAdmin)).Post("/", rt.orderHandler.Create)

			// Waiter can close orders
			r.With(middleware.RequireRole(domain.RoleWaiter, domain.RoleAdmin)).Put("/{id}/close", rt.orderHandler.CloseOrder)

			// Cook and Admin can update status
			r.With(middleware.RequireRole(domain.RoleCook, domain.RoleAdmin)).Put("/{id}/status", rt.orderHandler.UpdateStatus)
		})

		// Ingredient routes (Admin only)
		r.Route("/api/ingredients", func(r chi.Router) {
			r.Use(middleware.RequireRole(domain.RoleAdmin))
			r.Get("/", rt.ingredientHandler.GetAll)
			r.Get("/low-stock", rt.ingredientHandler.GetLowStock)
			r.Get("/{id}", rt.ingredientHandler.GetByID)
			r.Post("/", rt.ingredientHandler.Create)
			r.Put("/{id}", rt.ingredientHandler.Update)
			r.Delete("/{id}", rt.ingredientHandler.Delete)
		})

		// Supply routes (Admin only)
		r.Route("/api/supplies", func(r chi.Router) {
			r.Use(middleware.RequireRole(domain.RoleAdmin))
			r.Get("/", rt.supplyHandler.GetAll)
			r.Post("/", rt.supplyHandler.Create)
		})

		// Table routes
		r.Route("/api/tables", func(r chi.Router) {
			r.Get("/", rt.tableHandler.GetAll)
			r.Get("/{id}", rt.tableHandler.GetByID)

			// Waiter and Admin can update status
			r.With(middleware.RequireRole(domain.RoleWaiter, domain.RoleAdmin)).Put("/{id}/status", rt.tableHandler.UpdateStatus)

			// Admin only
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole(domain.RoleAdmin))
				r.Post("/", rt.tableHandler.Create)
				r.Delete("/{id}", rt.tableHandler.Delete)
			})
		})

		// Analytics routes (Admin and Manager only)
		r.Route("/api/analytics", func(r chi.Router) {
			r.Use(middleware.RequireRole(domain.RoleAdmin, domain.RoleManager))

			// Dashboard - main endpoint
			r.Get("/dashboard", rt.analyticsHandler.GetDashboard)

			// Sales analytics
			r.Get("/sales/summary", rt.analyticsHandler.GetSalesSummary)
			r.Get("/sales/by-category", rt.analyticsHandler.GetSalesByCategory)
			r.Get("/sales/hourly", rt.analyticsHandler.GetHourlyRevenue)

			// Dishes analytics
			r.Get("/dishes/popular", rt.analyticsHandler.GetPopularDishes)

			// Orders analytics
			r.Get("/orders/stats", rt.analyticsHandler.GetOrderStats)

			// Staff analytics
			r.Get("/waiters/performance", rt.analyticsHandler.GetWaiterPerformance)

			// Inventory analytics
			r.Get("/ingredients/turnover", rt.analyticsHandler.GetIngredientTurnover)

			// Tables analytics
			r.Get("/tables/utilization", rt.analyticsHandler.GetTableUtilization)
		})

	})

	return r
}
