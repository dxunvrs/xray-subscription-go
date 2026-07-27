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

func (h *AdminUserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, users, http.StatusOK)
}

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
