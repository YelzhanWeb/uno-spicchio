// // file: internal/controller/http/handler.go

package httpAdapter

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log/slog"
// 	"net/http"
// 	"strconv"

// 	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
// 	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
// 	"github.com/YelzhanWeb/uno-spicchio/internal/usecase"
// 	"github.com/go-chi/chi/v5"
// )

// // Handler - это контейнер для всех зависимостей слоя HTTP.
// type Handler struct {
// 	usecases     *usecase.UseCases
// 	categoryRepo ports.CategoryRepository // Вы правильно добавили это поле
// }

// // NewHandler - конструктор для наших обработчиков.
// // ИСПРАВЛЕНО: Теперь он принимает categoryRepo.
// func NewHandler(usecases *usecase.UseCases, categoryRepo ports.CategoryRepository) *Handler {
// 	return &Handler{
// 		usecases:     usecases,
// 		categoryRepo: categoryRepo, // И сохраняет его в структуру
// 	}
// }

// // PostUserHandler - оставляем ваш метод-заглушку.
// func (h *Handler) PostUserHandler() {
// 	// ...
// }

// // createOrder ... (без изменений)
// func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
// 	var req CreateOrderRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}
// 	order := domain.Order{
// 		WaiterID: req.WaiterID,
// 		TableID:  req.TableID,
// 		Notes:    req.Notes,
// 		Items:    make([]domain.OrderItem, len(req.Items)),
// 	}
// 	for i, item := range req.Items {
// 		order.Items[i] = domain.OrderItem{
// 			DishID: item.DishID,
// 			Qty:    item.Qty,
// 			Notes:  item.Notes,
// 		}
// 	}
// 	orderID, err := h.usecases.Order.CreateOrder(r.Context(), order)
// 	if err != nil {
// 		slog.Error("Failed to create order", "error", err)
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"message":  "Order created successfully",
// 		"order_id": orderID,
// 	})
// }

// // createDish ... (без изменений)
// func (h *Handler) createDish(w http.ResponseWriter, r *http.Request) {
// 	var dish domain.Dish
// 	if err := json.NewDecoder(r.Body).Decode(&dish); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}
// 	dishID, err := h.usecases.Dish.CreateDish(r.Context(), dish)
// 	if err != nil {
// 		slog.Error("Failed to create dish", "error", err)
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"message": "Dish created successfully",
// 		"dish_id": dishID,
// 	})
// }

// // getActiveOrders ... (без изменений)
// func (h *Handler) getActiveOrders(w http.ResponseWriter, r *http.Request) {
// 	orders, err := h.usecases.Order.GetActiveOrders(r.Context())
// 	if err != nil {
// 		slog.Error("Failed to get active orders", "error", err)
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(orders)
// }

// // updateOrderStatus ... (без изменений)
// func (h *Handler) updateOrderStatus(w http.ResponseWriter, r *http.Request) {
// 	orderIDStr := chi.URLParam(r, "id")
// 	orderID, err := strconv.Atoi(orderIDStr)
// 	if err != nil {
// 		http.Error(w, "Invalid order ID", http.StatusBadRequest)
// 		return
// 	}
// 	var req struct {
// 		Status string `json:"status"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}
// 	if err := h.usecases.Order.UpdateOrderStatus(r.Context(), orderID, req.Status); err != nil {
// 		slog.Error("Failed to update order status", "error", err)
// 		if err.Error() == fmt.Sprintf("order with id %d not found", orderID) {
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 		} else {
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		}
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// }

// // getMenu ... (без изменений)
// func (h *Handler) getMenu(w http.ResponseWriter, r *http.Request) {
// 	menu, err := h.categoryRepo.GetAllWithDishes(r.Context())
// 	if err != nil {
// 		slog.Error("Failed to get menu", "error", err)
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(menu)
// }
