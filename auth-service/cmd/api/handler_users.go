package main

import (
	"auth-service/cmd/api/data"
	"auth-service/internal/database"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type NoContent struct{}

func (app *Config) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	validation := &data.Validation{}

	type Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	body := Body{}
	err := decoder.Decode(&body)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	err = validation.ValidateEmail(body.Email)
	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Email invalid: %v", err))
		return
	}

	password, err := app.Models.User.Password.Set(body.Password)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	user, err := app.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:    body.Email,
		Password: password,
	})

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	responseJSON(w, 201, app.Models.ToUser(user))
}

func (app *Config) handlerGetUserById(w http.ResponseWriter, r *http.Request) {

	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error parsing parameters: %v", err))
		return
	}

	uuid, err := uuid.Parse(paths[len(paths)-1])

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error parsing parameters: %v", err))
		return
	}

	user, err := app.DB.GetUserById(r.Context(), uuid)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error getting user by id: %v", err))
		return
	}

	responseJSON(w, 200, app.Models.ToUser(user))
}

func (app *Config) handlerGetUserByEmail(w http.ResponseWriter, r *http.Request) {

	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error parsing parameters: %v", err))
		return
	}

	email := paths[len(paths)-1]

	user, err := app.DB.GetUserByEmail(r.Context(), email)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error getting user by email: %v", err))
		return
	}

	responseJSON(w, 200, app.Models.ToUser(user))
}

func (app *Config) handlerGetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := app.DB.GetAllUsers(r.Context())

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error getting all users: %v", err))
		return
	}

	responseJSON(w, 200, app.Models.ToUsers(users))
}

func (app *Config) handlerGetAllActiveUsers(w http.ResponseWriter, r *http.Request) {

	users, err := app.DB.GetAllActiveUsers(r.Context())

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error getting all active users: %v", err))
		return
	}

	responseJSON(w, 200, app.Models.ToUsers(users))
}

func (app *Config) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error parsing parameters: %v", err))
		return
	}

	uuid, err := uuid.Parse(paths[len(paths)-1])

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error parsing parameters: %v", err))
		return
	}

	err = app.DB.DeleteUser(r.Context(), uuid)
	if err != nil {
		errorJSON(w, 400, fmt.Sprintf("Error deleting user: %v", err))
		return
	}

	responseJSON(w, 204, NoContent{})
}

func (app *Config) paramParser(w http.ResponseWriter, r *http.Request) ([]string, error) {
	paths := strings.Split(r.URL.Path, "/")

	if len(paths) == 0 {
		return nil, errors.New("Missing parameters")
	}

	return paths, nil
}
