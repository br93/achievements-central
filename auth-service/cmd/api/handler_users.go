package main

import (
	"auth-service/cmd/api/data"
	"auth-service/internal/database"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type NoContent struct{}

const (
	errFormat          = "%s: %v"
	errJSON            = "Error parsing JSON"
	errParams          = "Error parsing parameters"
	errEmail           = "Email invalid"
	errEmailOrPassword = "Email or password invalid"
	errCreate          = "Error creating user"
	errUpdate          = "Error updating user"
	errGet             = "Error getting user"
	errDelete          = "Error deleting user"
)

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
		errorJSON(w, 400, fmt.Sprintf(errFormat, errJSON, err))
		return
	}

	err = validation.ValidateEmail(body.Email)
	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errEmail, err))
		return
	}

	password, err := app.Models.User.Password.Set(body.Password)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errCreate, err))
		return
	}

	user, err := app.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:    body.Email,
		Password: password,
	})

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errCreate, err))
		return
	}

	responseJSON(w, 201, app.Models.ToUser(user))
}

func (app *Config) handlerGetUserById(w http.ResponseWriter, r *http.Request) {

	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errJSON, err))
		return
	}

	uuid, err := uuid.Parse(paths[len(paths)-1])

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errJSON, err))
		return
	}

	user, err := app.DB.GetUserById(r.Context(), uuid)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errGet+" by id", err))
		return
	}

	responseJSON(w, 200, app.Models.ToUser(user))
}

func (app *Config) handlerGetUserByEmail(w http.ResponseWriter, r *http.Request) {

	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	email := paths[len(paths)-1]

	user, err := app.DB.GetUserByEmail(r.Context(), email)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errGet+" by email", err))
		return
	}

	responseJSON(w, 200, app.Models.ToUser(user))
}

func (app *Config) handlerGetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := app.DB.GetAllUsers(r.Context())

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errGet+" all", err))
		return
	}

	responseJSON(w, 200, app.Models.ToUsers(users))
}

func (app *Config) handlerGetAllActiveUsers(w http.ResponseWriter, r *http.Request) {

	users, err := app.DB.GetAllActiveUsers(r.Context())

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errGet+" all active", err))
		return
	}

	responseJSON(w, 200, app.Models.ToUsers(users))
}

func (app *Config) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	uuid, err := uuid.Parse(paths[len(paths)-1])

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	err = app.DB.DeleteUser(r.Context(), uuid)
	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errDelete, err))
		return
	}

	responseJSON(w, 204, NoContent{})
}

func (app *Config) handlerUpdateEmail(w http.ResponseWriter, r *http.Request) {

	type Body struct {
		NewEmail string `json:"new_email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	body := Body{}
	err := decoder.Decode(&body)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errJSON, err))
		return
	}

	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	uuid, err := uuid.Parse(paths[len(paths)-2])

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	user, err := app.DB.GetUserById(r.Context(), uuid)

	match, _ := app.Models.User.Password.Matches(body.Password, user.Password)
	if !match {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errEmailOrPassword, nil))
		return
	}

	err = app.DB.UpdateEmail(r.Context(), database.UpdateEmailParams{
		Email: body.NewEmail,
		ID:    uuid,
	})

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errUpdate, err))
		return
	}

	responseJSON(w, 204, NoContent{})
}

func (app *Config) handlerUpdatePassword(w http.ResponseWriter, r *http.Request) {

	type Body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	decoder := json.NewDecoder(r.Body)
	body := Body{}
	err := decoder.Decode(&body)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errJSON, err))
		return
	}

	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	uuid, err := uuid.Parse(paths[len(paths)-2])

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	user, err := app.DB.GetUserById(r.Context(), uuid)

	match, _ := app.Models.User.Password.Matches(body.OldPassword, user.Password)
	if !match {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errEmailOrPassword, nil))
		return
	}

	password, err := app.Models.User.Password.Set(body.NewPassword)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errEmailOrPassword, nil))
		return
	}

	err = app.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		Password: password,
		ID:       uuid,
	})

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errUpdate, err))
		return
	}

	responseJSON(w, 204, NoContent{})
}

func (app *Config) handlerUpdateSuperUser(w http.ResponseWriter, r *http.Request) {
	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	uuid, err := uuid.Parse(paths[len(paths)-2])

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	err = app.DB.UpdateSuperUser(r.Context(), database.UpdateSuperUserParams{
		IsSuperuser: sql.NullBool{Bool: true, Valid: true},
		ID:          uuid,
	})

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errUpdate, err))
		return
	}

	responseJSON(w, 204, NoContent{})
}

func (app *Config) handlerUpdateActive(w http.ResponseWriter, r *http.Request) {
	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	uuid, err := uuid.Parse(paths[len(paths)-2])

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	err = app.DB.UpdateActive(r.Context(), database.UpdateActiveParams{
		IsActive: sql.NullBool{Bool: true, Valid: true},
		ID:       uuid,
	})

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errUpdate, err))
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
