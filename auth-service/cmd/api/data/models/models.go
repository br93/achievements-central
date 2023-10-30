package data

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Models struct {
	User User
}

type User struct {
	ID          uint         `json:"id"`
	UserID      string       `json:"user_id"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Password    string       `json:"password"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	ActivatedAt sql.NullTime `json:"activated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at,omitempty"`
}

func (this *User) HashPassword() ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(this.Password), 12)
	return hashedPassword, err
}

func (this *User) CheckHash(text string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(this.Password), []byte(text))

	if err != nil {
		return false, err
	}

	return true, nil
}
