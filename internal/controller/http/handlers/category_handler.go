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

type CategoryHandler struct {
	categoryService ports.CategoryService
}

func NewCategoryHandler(categoryService ports.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.GetAll(r.Context())
	if err != nil {
		response.InternalError(w, "failed to get categories")
		return
	}

	response.Success(w, categories)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid category id")
		return
	}

	category, err := h.categoryService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrCategoryNotFound {
			response.NotFound(w, "category not found")
			return
		}
		response.InternalError(w, "failed to get category")
		return
	}

	response.Success(w, category)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.categoryService.Create(r.Context(), &category); err != nil {
		response.InternalError(w, "failed to create category")
		return
	}

	response.Created(w, category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid category id")
		return
	}

	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	category.ID = id
	if err := h.categoryService.Update(r.Context(), &category); err != nil {
		if err == domain.ErrCategoryNotFound {
			response.NotFound(w, "category not found")
			return
		}
		response.InternalError(w, "failed to update category")
		return
	}

	response.Success(w, category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid category id")
		return
	}

	if err := h.categoryService.Delete(r.Context(), id); err != nil {
		if err == domain.ErrCategoryNotFound {
			response.NotFound(w, "category not found")
			return
		}
		response.InternalError(w, "failed to delete category")
		return
	}

	response.Success(w, map[string]string{"message": "category deleted"})
}
