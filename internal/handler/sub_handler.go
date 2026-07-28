package handler

import (
	"errors"
	"net/http"
	"xray-subscription-go/internal/service"

	"github.com/go-chi/chi/v5"
)

type SubHandler struct {
	subService *service.SubscriptionService
}

func NewSubHandler(subService *service.SubscriptionService) *SubHandler {
	return &SubHandler{
		subService: subService,
	}
}

// GetSubscription godoc
// @Summary      Получить подписку
// @Tags         Subscription
// @Produce      text/plain
// @Param        email path string true "Email пользователя"
// @Success      200 {string} string "Vless ссылка"
// @Router       /sub/{email} [get]
func (h *SubHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		http.Error(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	content, err := h.subService.GetSubscription(email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "Subscription not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(content))
}
