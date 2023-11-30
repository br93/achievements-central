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
	errJSON            = "error parsing JSON "
	errParams          = "error parsing parameters "
	errEmail           = "email invalid "
	errEmailOrPassword = "email or password invalid "
	errCreate          = "error creating user "
	errUpdate          = "error updating user "
	errGet             = "error getting user "
	errDelete          = "error deleting user "
	errLogin           = "error login "
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) handlerLogin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	body := LoginBody{}
	err := decoder.Decode(&body)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errJSON, err))
		return
	}

	user, err := app.Service.Login(r.Context(), body)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errEmailOrPassword))
		return
	}

	token, err := app.GenerateToken(user.Email)
	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, "failed to create token ", err))
	}

	app.setCookie(w, "jwt", token, 3600, false, true)
	responseJSON(w, 200, user)
}

func (app *Config) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	validation := &data.Validation{}

	decoder := json.NewDecoder(r.Body)
	body := LoginBody{}
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

	password, err := app.Service.SetPassword(body.Password)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errCreate, err))
		return
	}

	user, err := app.Service.Create(r.Context(), database.CreateUserParams{
		Email:    body.Email,
		Password: password,
	})

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errCreate, err))
		return
	}

	responseJSON(w, 201, user)
}

func (app *Config) handleGetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")

	if user == nil {
		errorJSON(w, 401, fmt.Sprintf("Unauthorized"))
		return
	}

	responseJSON(w, 200, user)
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

	user, err := app.Service.GetById(r.Context(), uuid)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errGet+" by id", err))
		return
	}

	responseJSON(w, 200, user)
}

func (app *Config) handlerGetUserByEmail(w http.ResponseWriter, r *http.Request) {

	paths, err := app.paramParser(w, r)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errParams, err))
		return
	}

	email := paths[len(paths)-1]

	user, err := app.Service.GetByEmail(r.Context(), email)

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errGet+" by email", err))
		return
	}

	responseJSON(w, 200, user)
}

func (app *Config) handlerGetAllUsers(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value("user").(data.User)
	fmt.Println(u)
	users, err := app.Service.GetAll(r.Context())

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errGet+" all", err))
		return
	}

	responseJSON(w, 200, users)
}

func (app *Config) handlerGetAllActiveUsers(w http.ResponseWriter, r *http.Request) {

	users, err := app.Service.GetAllActive(r.Context())

	if err != nil {
		errorJSON(w, 400, fmt.Sprintf(errFormat, errGet+" all active", err))
		return
	}

	responseJSON(w, 200, users)
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

	err = app.Service.Delete(r.Context(), uuid)
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

	err = app.Service.UpdateEmail(r.Context(), body.Password, database.UpdateEmailParams{
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

	err = app.Service.UpdatePassword(r.Context(), body.OldPassword, database.UpdatePasswordParams{
		Password: body.NewPassword,
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

	err = app.Service.UpdateSuperUser(r.Context(), database.UpdateSuperUserParams{
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

	err = app.Service.UpdateActiveUser(r.Context(), database.UpdateActiveParams{
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
