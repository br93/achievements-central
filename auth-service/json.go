package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseJSON(w http.ResponseWriter, code int, payload interface{}) {

	type response struct {
		Success string      `json:"success"`
		Data    interface{} `json:"data"`
	}

	data, err := json.Marshal(response{
		Success: "true",
		Data:    payload,
	})

	if err != nil {
		log.Printf("Marshal failed: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func errorJSON(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Printf("Server error: %v", message)
	}

	type error struct {
		Success string `json:"success"`
		Error   string `json:"error"`
	}

	responseJSON(w, code, error{
		Success: "false",
		Error:   message,
	})
}
