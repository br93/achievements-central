package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	app.v1Routes(router)

	return router
}

func (app *Config) v1Routes(router *chi.Mux) {
	router.Get("/health", app.handlerHealth)
	app.userRoutes(router)
	router.Mount("/api/v1", router)
}

func (app *Config) userRoutes(router *chi.Mux) {
	router.Get("/users", app.handlerGetAllUsers)
	router.Get("/users/{id}", app.handlerGetUserById)
	router.Get("/users/active", app.handlerGetAllActiveUsers)
	router.Post("/users", app.handlerCreateUser)
	router.Put("/users/{id}/email", app.handlerUpdateEmail)
	router.Put("/users/{id}/password", app.handlerUpdatePassword)
	router.Put("/users/{id}/active", app.handlerUpdateActive)
	router.Put("/users/{id}/super-user", app.handlerUpdateSuperUser)
	router.Delete("/users/{id}", app.handlerDeleteUser)
}
