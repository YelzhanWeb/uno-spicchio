// file: internal/controller/http/router.go

package http

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter создает новый роутер chi и настраивает все маршруты.
func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	// Используем стандартные middleware для логирования запросов,
	// восстановления после паник и т.д.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes) // Убирает слэш в конце URL

	// Группа маршрутов для нашего API v1
	r.Route("/api/v1", func(r chi.Router) {
		// --- Маршруты для заказов ---
		r.Post("/orders", h.createOrder) // НАШ НОВЫЙ МАРШРУТ!
		r.Get("/orders", h.getActiveOrders)
		r.Patch("/orders/{id}/status", h.updateOrderStatus) // <-- ДОБАВЬТЕ ЭТУ СТРОКУ
		r.Get("/menu", h.getMenu)
		// --- Маршруты для блюд ---
		r.Post("/dishes", h.createDish)
		// В будущем здесь будут другие маршруты:
		// r.Get("/orders", h.getOrders)
		// r.Get("/orders/{orderID}", h.getOrderById)

		// --- Маршруты для пользователей ---
		// r.Post("/users", h.createUser)
		// r.Get("/users", h.getUsers)
	})

	// Настраиваем раздачу статических файлов (HTML, CSS, JS) из папки 'static'
	staticPath, _ := filepath.Abs("./static/")
	fs := http.FileServer(http.Dir(staticPath))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Перенаправляем с главной страницы на страницу заказов
	// Перенаправляем с главной страницы на страницу СОЗДАНИЯ заказа
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/create-order.html", http.StatusMovedPermanently)
	})
	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// Этот хелпер взят из документации chi для корректной работы с файлами.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
