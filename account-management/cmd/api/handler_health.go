package main

import (
	"net/http"
)

func (*Config) handlerHealth(w http.ResponseWriter, r *http.Request) {

	type Health struct {
		Status string `json:"status" `
	}

	responseJSON(w, 200, Health{
		Status: "UP",
	})
}
