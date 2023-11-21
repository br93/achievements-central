package data

import (
	"auth-service/internal/database"
	"time"

	"github.com/google/uuid"
)

type Models struct {
	User User
}

type User struct {
	ID        uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

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
