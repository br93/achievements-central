package main

import (
	"net/http"
	"net/url"

	"github.com/go-chi/jwtauth/v5"
)

func (*Config) setCookie(w http.ResponseWriter, name string, value string, maxAge int, secure bool, httpOnly bool) {

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     "/",
		SameSite: http.SameSite(maxAge),
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (*Config) getCookie(r *http.Request, name string) string {
	return jwtauth.TokenFromCookie(r)
}
