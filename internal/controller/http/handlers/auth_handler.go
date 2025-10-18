package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YelzhanWeb/uno-spicchio/internal/controller/http/middleware"
	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
	"github.com/YelzhanWeb/uno-spicchio/pkg/response"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	token, user, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials || err == domain.ErrUserNotActive {
			response.Unauthorized(w, err.Error())
			return
		}
		response.InternalError(w, "failed to login")
		return
	}
	fmt.Println(token, user)
	response.Success(w, LoginResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Unauthorized(w, "user not authenticated")
		return
	}

	user, err := h.authService.GetCurrentUser(r.Context(), userID)
	if err != nil {
		response.InternalError(w, "failed to get user")
		return
	}

	response.Success(w, user)
}
