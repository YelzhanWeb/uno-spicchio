package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/YelzhanWeb/uno-spicchio/internal/controller/http/middleware"
	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/internal/usecase"
	"github.com/YelzhanWeb/uno-spicchio/pkg/response"
	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	orderService ports.OrderService
}

func NewOrderHandler(orderService ports.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type CreateOrderRequest struct {
	TableNumber int                      `json:"table_number"`
	Notes       *string                  `json:"notes"`
	Items       []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	DishID int     `json:"dish_id"`
	Qty    int     `json:"qty"`
	Notes  *string `json:"notes"`
}

func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	statusStr := r.URL.Query().Get("status")
	var status *domain.OrderStatus
	if statusStr != "" {
		s := domain.OrderStatus(statusStr)
		status = &s
	}

	orders, err := h.orderService.GetAll(r.Context(), status)
	if err != nil {
		response.InternalError(w, "failed to get orders")
		return
	}

	response.Success(w, orders)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid order id")
		return
	}

	order, err := h.orderService.GetByID(r.Context(), id)
	if err != nil {
		if err == usecase.ErrOrderNotFound {
			response.NotFound(w, "order not found")
			return
		}
		response.InternalError(w, "failed to get order")
		return
	}

	response.Success(w, order)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	waiterID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Unauthorized(w, "user not authenticated")
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	order := &domain.Order{
		WaiterID:    waiterID,
		TableNumber: req.TableNumber,
		Notes:       req.Notes,
	}

	var items []domain.OrderItem
	for _, itemReq := range req.Items {
		items = append(items, domain.OrderItem{
			DishID: itemReq.DishID,
			Qty:    itemReq.Qty,
			Notes:  itemReq.Notes,
		})
	}

	if err := h.orderService.Create(r.Context(), order, items); err != nil {
		if err == usecase.ErrInsufficientStock {
			response.BadRequest(w, "insufficient stock for order")
			return
		}
		response.InternalError(w, "failed to create order")
		return
	}

	response.Created(w, order)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid order id")
		return
	}

	var req struct {
		Status domain.OrderStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.orderService.UpdateStatus(r.Context(), id, req.Status); err != nil {
		if err == usecase.ErrOrderNotFound {
			response.NotFound(w, "order not found")
			return
		}
		if err == usecase.ErrInvalidStatusChange {
			response.BadRequest(w, "invalid status change")
			return
		}
		response.InternalError(w, "failed to update order status")
		return
	}

	response.Success(w, map[string]string{"message": "order status updated"})
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid order id")
		return
	}

	if err := h.orderService.CloseOrder(r.Context(), id); err != nil {
		if err == usecase.ErrOrderNotFound {
			response.NotFound(w, "order not found")
			return
		}
		response.InternalError(w, "failed to close order")
		return
	}

	response.Success(w, map[string]string{"message": "order closed"})
}
