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

type IngredientHandler struct {
	ingredientService ports.IngredientService
}

func NewIngredientHandler(ingredientService ports.IngredientService) *IngredientHandler {
	return &IngredientHandler{ingredientService: ingredientService}
}

func (h *IngredientHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ingredients, err := h.ingredientService.GetAll(r.Context())
	if err != nil {
		response.InternalError(w, "failed to get ingredients")
		return
	}

	response.Success(w, ingredients)
}

func (h *IngredientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid ingredient id")
		return
	}

	ingredient, err := h.ingredientService.GetByID(r.Context(), id)
	if err != nil {
		response.InternalError(w, "failed to get ingredient")
		return
	}

	if ingredient == nil {
		response.NotFound(w, "ingredient not found")
		return
	}

	response.Success(w, ingredient)
}

func (h *IngredientHandler) GetLowStock(w http.ResponseWriter, r *http.Request) {
	ingredients, err := h.ingredientService.GetLowStock(r.Context())
	if err != nil {
		response.InternalError(w, "failed to get low stock ingredients")
		return
	}

	response.Success(w, ingredients)
}

func (h *IngredientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var ingredient domain.Ingredient
	if err := json.NewDecoder(r.Body).Decode(&ingredient); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.ingredientService.Create(r.Context(), &ingredient); err != nil {
		response.InternalError(w, "failed to create ingredient")
		return
	}

	response.Created(w, ingredient)
}

func (h *IngredientHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid ingredient id")
		return
	}

	var ingredient domain.Ingredient
	if err := json.NewDecoder(r.Body).Decode(&ingredient); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	ingredient.ID = id
	if err := h.ingredientService.Update(r.Context(), &ingredient); err != nil {
		response.InternalError(w, "failed to update ingredient")
		return
	}

	response.Success(w, ingredient)
}

func (h *IngredientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid ingredient id")
		return
	}

	if err := h.ingredientService.Delete(r.Context(), id); err != nil {
		response.InternalError(w, "failed to delete ingredient")
		return
	}

	response.Success(w, map[string]string{"message": "ingredient deleted"})
}
