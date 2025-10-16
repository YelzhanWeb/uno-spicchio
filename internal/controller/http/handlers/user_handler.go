package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/internal/usecase"
	"github.com/YelzhanWeb/uno-spicchio/pkg/response"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type CreateUserRequest struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Role     domain.Role `json:"role"`
	PhotoKey string      `json:"photo_key"`
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll(r.Context())
	if err != nil {
		response.InternalError(w, "failed to get users")
		return
	}

	response.Success(w, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid user id")
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		if err == usecase.ErrUserNotFound {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalError(w, "failed to get user")
		return
	}

	response.Success(w, user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	user := &domain.User{
		Username: req.Username,
		Role:     req.Role,
		PhotoKey: req.PhotoKey,
	}

	if err := h.userService.Create(r.Context(), user, req.Password); err != nil {
		if err == usecase.ErrUserExists {
			response.BadRequest(w, "user already exists")
			return
		}
		response.InternalError(w, "failed to create user")
		return
	}

	response.Created(w, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid user id")
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	user.ID = id
	if err := h.userService.Update(r.Context(), &user); err != nil {
		if err == usecase.ErrUserNotFound {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalError(w, "failed to update user")
		return
	}

	response.Success(w, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid user id")
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		if err == usecase.ErrUserNotFound {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalError(w, "failed to delete user")
		return
	}

	response.Success(w, map[string]string{"message": "user deleted"})
}
