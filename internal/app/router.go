package app

import (
	"net/http"

	"xray-subscription-go/internal/config"
	"xray-subscription-go/internal/handler"

	_ "xray-subscription-go/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(adminHandler *handler.AdminUserHandler, subHandler *handler.SubHandler, config *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Get("/sub/{email}", subHandler.GetSubscription)

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
