package data

import (
	"auth-service/internal/database"
	"log"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

type Models struct {
	User     User
	password password
}

type User struct {
	ID        uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type password struct{}

type LoggedUser struct {
	Email    string    `json:"email"`
	LoggedAt time.Time `json:"logged_at"`
}

func (*Models) ToUser(tbUser database.TbUser) User {
	return User{
		ID:        tbUser.ID,
		Email:     tbUser.Email,
		CreatedAt: tbUser.CreatedAt,
	}
}

func (*Models) ToLoggedUser(tbUser database.TbUser) LoggedUser {
	return LoggedUser{
		Email:    tbUser.Email,
		LoggedAt: time.Now(),
	}
}

func (model *Models) ToUsers(tbUsers []database.TbUser) []User {
	var users []User

	for _, data := range tbUsers {
		users = append(users, model.ToUser(data))
	}

	return users
}

func (*password) Set(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (*password) Matches(plainText string, password string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plainText, password)
	if err != nil {
		log.Fatal(err)
	}

	return match, nil
}
