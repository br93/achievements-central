package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

func responseJSON(w http.ResponseWriter, code int, payload interface{}) {

	type response struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}

	responseType := reflect.TypeOf(payload)

	data, err := json.Marshal(response{
		Success: responseType.Name() != "error",
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
		Status int    `json:"status"`
		Error  string `json:"error"`
	}

	responseJSON(w, code, error{
		Status: code,
		Error:  message,
	})
}
