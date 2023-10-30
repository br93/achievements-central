package repository

import (
	models "auth-service/cmd/api/data/models"
	"context"
	"database/sql"
	"errors"
	"log"
	"os/exec"
	"time"
)

type Models = models.Models
type User = models.User

const timeout = time.Second * 3
const failedScanMessage = "scanning failed"

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

func InsertUser(user User) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	userId, err := exec.Command("uuidgen").Output()
	if err != nil {
		return 0, errors.New("UUID generation failed")
	}

	var newId uint
	query := `insert into tb_user (user_id, username, email, password, created_at, updated_at) 
		values ($1, $2, $3, $4, $5, $6) return id`

	err = db.QueryRowContext(ctx, query,
		userId, user.Username, user.Email, user.Password, time.Now(), time.Now()).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

func GetAllUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `select * from tb_users 
		where activated_at is not null`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.UserID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.ActivatedAt,
			&user.DeletedAt)

		if err != nil {
			log.Println(failedScanMessage, err)
			return nil, err
		}

		users = append(users, &user)

	}

	return users, nil

}

func GetUserByUserId(userId string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `select * from tb_users 
		where user_id = $1 and activated_at is not null`

	var user User
	row := db.QueryRowContext(ctx, query, userId)

	err := row.Scan(
		&user.ID,
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.ActivatedAt,
		&user.DeletedAt)

	if err != nil {
		log.Println(failedScanMessage, err)
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `select * from tb_users 
		where username = $1 and activated_at is not null`

	var user User
	row := db.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.ActivatedAt,
		&user.DeletedAt)

	if err != nil {
		log.Println(failedScanMessage, err)
		return nil, err
	}

	return &user, nil
}

func UpdateUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `update tb_users set email = $1, password = $2, updated_at = $3, 
		where user_id = $4 and activated_at is not null`

	_, err := db.ExecContext(ctx, query,
		user.Email,
		user.Password,
		time.Now(),
		user.UserID)

	if err != nil {
		return err
	}

	return nil
}

func DeleteUserById(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `update tb_users set activated_at = $1, deleted_at = $2, updated_at = $3,
		where user_id = $3 and activated_at is not null`

	_, err := db.ExecContext(ctx, query,
		nil,
		time.Now(),
		time.Now(),
		userId)

	if err != nil {
		return err
	}

	return nil
}

func ActivateUserById(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `update tb_users set activated_at = $1, updated_at = $2,
		where user_id = $2 and activated_at is not null`

	_, err := db.ExecContext(ctx, query,
		time.Now(),
		time.Now(),
		userId)

	if err != nil {
		return err
	}

	return nil
}
