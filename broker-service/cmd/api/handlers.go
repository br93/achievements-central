package main

import (
	"net/http"
)

func (app *Config) Broker(writer http.ResponseWriter, request *http.Request) {
	response := Response{
		Success: true,
		Message: "Broker acknowledged",
	}

	app.writeJSON(writer, http.StatusOK, response)
}
