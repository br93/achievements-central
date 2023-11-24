package main

import (
	"context"
	"net/http"
)

func (app *Config) UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie := app.getCookie(r, "jwt")
		if cookie == "" {
			ctx := context.WithValue(r.Context(), "user", nil)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		token, err := app.Token.Decode(cookie)

		if err != nil {
			ctx := context.WithValue(r.Context(), "user", nil)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		email, _ := token.Get("user_email")
		user, err := app.Service.GetByEmail(r.Context(), email.(string))

		if err != nil {
			ctx := context.WithValue(r.Context(), "user", nil)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
