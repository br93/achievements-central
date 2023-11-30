package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Cookie", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	app.v1Routes(router)

	return router
}

func (app *Config) v1Routes(router *chi.Mux) {

	app.publicRoutes(router)

	router.Mount("/api/v1", router)
}

func (app *Config) publicRoutes(router *chi.Mux) {
	router.Group(func(r chi.Router) {
		r.Get("/health", app.handlerHealth)
	})
}
