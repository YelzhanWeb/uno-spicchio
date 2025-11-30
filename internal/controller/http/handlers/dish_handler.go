package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/pkg/response"
	"github.com/go-chi/chi/v5"
)

type DishHandler struct {
	dishService ports.DishService
}

func NewDishHandler(dishService ports.DishService) *DishHandler {
	return &DishHandler{dishService: dishService}
}

func (h *DishHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"

	dishes, err := h.dishService.GetAll(r.Context(), activeOnly)
	if err != nil {
		response.InternalError(w, "failed to get dishes")
		return
	}

	response.Success(w, dishes)
}

func (h *DishHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid dish id")
		return
	}

	dish, err := h.dishService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrDishNotFound {
			response.NotFound(w, "dish not found")
			return
		}
		response.InternalError(w, "failed to get dish")
		return
	}

	response.Success(w, dish)
}

func (h *DishHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dish domain.Dish
	if err := json.NewDecoder(r.Body).Decode(&dish); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.dishService.Create(r.Context(), &dish); err != nil {
		response.InternalError(w, "failed to create dish")
		return
	}

	response.Created(w, dish)
}

func (h *DishHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid dish id")
		return
	}

	var dish domain.Dish
	if err := json.NewDecoder(r.Body).Decode(&dish); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	dish.ID = id
	if err := h.dishService.Update(r.Context(), &dish); err != nil {
		if err == domain.ErrDishNotFound {
			response.NotFound(w, "dish not found")
			return
		}
		response.InternalError(w, "failed to update dish")
		return
	}

	response.Success(w, dish)
}
func (h *DishHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid dish id")
		return
	}

	if err := h.dishService.Delete(r.Context(), id); err != nil {
		if err == domain.ErrDishNotFound {
			response.NotFound(w, "dish not found")
			return
		}
		response.InternalError(w, "failed to delete dish")
		return
	}

	response.Success(w, map[string]string{"message": "dish deleted"})
}

func (h *DishHandler) GetIngredients(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid dish id")
		return
	}

	ingredients, err := h.dishService.GetIngredients(r.Context(), id)
	if err != nil {
		response.InternalError(w, "failed to get dish ingredients")
		return
	}

	response.Success(w, ingredients)
}
