package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"xray-subscription-go/internal/dto"
	"xray-subscription-go/internal/service"

	"github.com/go-chi/chi/v5"
)

type AdminUserHandler struct {
	userService *service.UserService
}

func NewAdminUserHandler(userService *service.UserService) *AdminUserHandler {
	return &AdminUserHandler{userService: userService}
}

// CreateUser godoc
// @Summary      Создать нового пользователя
// @Description  Генерирует новый UUID и создает пользователя в Xray и базе данных
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BasicAuth
// @Param        request body dto.CreateUserRequest true "Данные для создания пользователя"
// @Success      201 {object} dto.UserResponse
// @Router       /api/admin/users [post]
func (h *AdminUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, `{"error": "Email is required"}`, http.StatusBadRequest)
		return
	}

	createdUser, err := h.userService.CreateUser(req.Email)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			writeJSONError(w, "User already exists", http.StatusConflict)
			return
		}
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, createdUser, http.StatusCreated)
}

// GetAllUsers godoc
// @Summary      Получить список всех пользователей
// @Description  Возвращает массив всех зарегистрированных пользователей
// @Tags         Admin
// @Produce      json
// @Security     BasicAuth
// @Success      200 {array} dto.UserResponse
// @Router       /api/admin/users [get]
func (h *AdminUserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, users, http.StatusOK)
}

// GetUserByEmail godoc
// @Summary      Получить пользователя по Email
// @Description  Возвращает данные конкретного пользователя по его email
// @Tags         Admin
// @Produce      json
// @Security     BasicAuth
// @Param        email path string true "Email пользователя"
// @Success      200 {object} dto.UserResponse
// @Router       /api/admin/users/{email} [get]
func (h *AdminUserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		writeJSONError(w, "Email is required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			writeJSONError(w, "User not found", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, user, http.StatusOK)
}

// DeleteUserByEmail godoc
// @Summary      Удалить пользователя по Email
// @Description  Удаляет пользователя из Xray и базы данных
// @Tags         Admin
// @Security     BasicAuth
// @Param        email path string true "Email пользователя"
// @Success      204 "No Content"
// @Router       /api/admin/users/{email} [delete]
func (h *AdminUserHandler) DeleteUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		writeJSONError(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	err := h.userService.DeleteUser(email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			writeJSONError(w, "User not found", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	writeJSON(w, map[string]string{"error": message}, statusCode)
}
