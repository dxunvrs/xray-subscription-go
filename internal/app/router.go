package app

import (
	"net/http"

	"xray-subscription-go/internal/config"
	"xray-subscription-go/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(adminHandler *handler.AdminUserHandler, config *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/admin/users", func(r chi.Router) {
		r.Use(middleware.BasicAuth("Admin Area", map[string]string{
			config.AdminUsername: config.AdminPassword,
		}))

		r.Post("/", adminHandler.CreateUser)
		r.Get("/", adminHandler.GetAllUsers)
		r.Get("/{email}", adminHandler.GetUserByEmail)
		r.Delete("/{email}", adminHandler.DeleteUserByEmail)
	})

	return r
}
