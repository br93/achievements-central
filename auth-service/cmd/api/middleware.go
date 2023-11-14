package main

import (
	"context"
	"net/http"
	"strings"
)

func (app *Config) UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := strings.Split(r.Header.Get("Authorization"), " ")
		token, err := app.Token.Decode(auth[1])

		if err != nil {
			ctx := context.WithValue(r.Context(), "user", nil)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		email, _ := token.Get("user_email")
		user, err := app.DB.GetUserByEmail(r.Context(), email.(string))

		if err != nil {
			ctx := context.WithValue(r.Context(), "user", nil)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		ctx := context.WithValue(r.Context(), "user", app.Models.ToUser(user))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
