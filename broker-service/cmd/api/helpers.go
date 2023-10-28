package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload,omitempty"`
}

const headerKey = "Content-Type"
const headerValue = "application/json"

func (app *Config) readJSON(writer http.ResponseWriter, request *http.Request, data interface{}) error {

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("JSON invalid")
	}

	return nil
}

func (app *Config) writeJSON(writer http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	output, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return app.setHeaders(writer, status, output, headers)
}

func (app *Config) errorJSON(writer http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	response := Response{
		Success: false,
		Message: err.Error(),
	}

	return app.writeJSON(writer, statusCode, response)
}

func (app *Config) setHeaders(writer http.ResponseWriter, status int, output []byte, headers []http.Header) error {
	if headers != nil && len(headers) > 0 {
		for key, value := range headers[0] {
			writer.Header()[key] = value
		}
	}

	writer.Header().Set(headerKey, headerValue)
	writer.WriteHeader(status)
	_, err := writer.Write(output)
	if err != nil {
		return err
	}

	return nil
}
