package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
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
		response.InternalError(w, "failed to get users: "+err.Error())
		return
	}

	response.Success(w, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid user id: must be a number")
		return
	}

	if id <= 0 {
		response.BadRequest(w, "user id must be positive")
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalError(w, "failed to get user: "+err.Error())
		return
	}

	response.Success(w, user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body: "+err.Error())
		return
	}

	// Валидация
	if strings.TrimSpace(req.Username) == "" {
		response.BadRequest(w, "username is required")
		return
	}

	if len(req.Username) < 3 {
		response.BadRequest(w, "username must be at least 3 characters")
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		response.BadRequest(w, "password is required")
		return
	}

	if len(req.Password) < 6 {
		response.BadRequest(w, "password must be at least 6 characters")
		return
	}

	validRoles := map[domain.Role]bool{
		domain.RoleAdmin:   true,
		domain.RoleManager: true,
		domain.RoleWaiter:  true,
		domain.RoleCook:    true,
	}

	if !validRoles[req.Role] {
		response.BadRequest(w, "invalid role. Must be: admin, manager, waiter, or cook")
		return
	}

	user := &domain.User{
		Username: strings.TrimSpace(req.Username),
		Role:     req.Role,
		PhotoKey: strings.TrimSpace(req.PhotoKey),
	}

	if err := h.userService.Create(r.Context(), user, req.Password); err != nil {
		if err == domain.ErrUserExists {
			response.BadRequest(w, "user with this username already exists")
			return
		}
		response.InternalError(w, "failed to create user: "+err.Error())
		return
	}

	response.Created(w, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid user id: must be a number")
		return
	}

	if id <= 0 {
		response.BadRequest(w, "user id must be positive")
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.BadRequest(w, "invalid request body: "+err.Error())
		return
	}

	// Валидация
	if strings.TrimSpace(user.Username) == "" {
		response.BadRequest(w, "username is required")
		return
	}

	if len(user.Username) < 3 {
		response.BadRequest(w, "username must be at least 3 characters")
		return
	}

	validRoles := map[domain.Role]bool{
		domain.RoleAdmin:   true,
		domain.RoleManager: true,
		domain.RoleWaiter:  true,
		domain.RoleCook:    true,
	}

	if !validRoles[user.Role] {
		response.BadRequest(w, "invalid role. Must be: admin, manager, waiter, or cook")
		return
	}

	user.ID = id
	user.Username = strings.TrimSpace(user.Username)
	user.PhotoKey = strings.TrimSpace(user.PhotoKey)

	if err := h.userService.Update(r.Context(), &user); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalError(w, "failed to update user: "+err.Error())
		return
	}

	response.Success(w, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid user id: must be a number")
		return
	}

	if id <= 0 {
		response.BadRequest(w, "user id must be positive")
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalError(w, "failed to delete user: "+err.Error())
		return
	}

	response.Success(w, map[string]string{"message": "user deleted successfully"})
}
