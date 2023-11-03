package main

import (
	"auth-service/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *Config) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type Parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	password, err := app.Models.User.Password.Set(params.Password)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	user, err := app.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:    params.Email,
		Password: password,
	})

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	responseJSON(w, 201, app.Models.ToUser(user))
}
