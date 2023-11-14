package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
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

	app.publicRoutes(router)
	app.privateRoutes(router)

	router.Mount("/api/v1", router)
}

func (app *Config) publicRoutes(router *chi.Mux) {
	router.Group(func(r chi.Router) {
		r.Get("/health", app.handlerHealth)
		r.Post("/users", app.handlerCreateUser)
		r.Post("/login", app.handlerLogin)
	})

}

func (app *Config) privateRoutes(router *chi.Mux) {
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(app.Token))
		r.Use(jwtauth.Authenticator)
		r.Use(app.UserContext)

		r.Get("/users", app.handlerGetAllUsers)
		r.Get("/users/{id}", app.handlerGetUserById)
		r.Get("/users/active", app.handlerGetAllActiveUsers)
		r.Get("/users/user", app.handleGetUser)

		r.Put("/users/{id}/email", app.handlerUpdateEmail)
		r.Put("/users/{id}/password", app.handlerUpdatePassword)
		r.Put("/users/{id}/active", app.handlerUpdateActive)
		r.Put("/users/{id}/super-user", app.handlerUpdateSuperUser)

		r.Delete("/users/{id}", app.handlerDeleteUser)
	})

}
