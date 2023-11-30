package main

import (
	"auth-service/cmd/api/data"
	"auth-service/internal/database"
	"context"
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

var passwordError = "password doesn't match"

type Service struct {
	model      data.Models
	repository *database.Queries
}

func NewService(rep *database.Queries) *Service {
	return &Service{
		repository: rep,
	}
}

func (s *Service) Login(ctx context.Context, params LoginBody) (*data.LoggedUser, error) {
	u, err := s.repository.GetUserByEmail(ctx, params.Email)

	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	match, _ := s.MatchPassword(params.Password, u.Password)
	if !match {
		return nil, fmt.Errorf(passwordError, err)
	}

	user := s.model.ToLoggedUser(u)

	return &user, nil
}

func (s *Service) Create(ctx context.Context, params database.CreateUserParams) (*data.User, error) {
	response, err := s.repository.CreateUser(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	user := s.model.ToUser(response)
	return &user, nil
}

func (s *Service) GetAll(ctx context.Context) (*[]data.User, error) {
	response, err := s.repository.GetAllUsers(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}

	users := s.model.ToUsers(response)
	return &users, nil
}

func (s *Service) GetAllActive(ctx context.Context) (*[]data.User, error) {
	response, err := s.repository.GetAllActiveUsers(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all active users: %w", err)
	}

	users := s.model.ToUsers(response)
	return &users, nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (*data.User, error) {
	response, err := s.repository.GetUserById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}

	user := s.model.ToUser(response)
	return &user, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*data.User, error) {
	response, err := s.repository.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	user := s.model.ToUser(response)
	return &user, nil
}

func (s *Service) UpdateEmail(ctx context.Context, password string, params database.UpdateEmailParams) error {

	user, err := s.repository.GetUserByEmail(ctx, params.Email)

	if err != nil {
		return fmt.Errorf("error getting user by email: %w", err)
	}

	match, _ := s.MatchPassword(password, user.Password)
	if !match {
		return fmt.Errorf(passwordError, err)
	}

	err = s.repository.UpdateEmail(ctx, params)

	if err != nil {
		return fmt.Errorf("error updating email: %w", err)
	}

	return nil
}

func (s *Service) UpdatePassword(ctx context.Context, oldPassword string, params database.UpdatePasswordParams) error {

	user, err := s.repository.GetUserById(ctx, params.ID)

	if err != nil {
		return fmt.Errorf("error getting user by id: %w", err)
	}

	match, _ := s.MatchPassword(oldPassword, user.Password)
	if !match {
		return fmt.Errorf("password doesn't match: %w", err)
	}

	password, err := s.SetPassword(params.Password)
	if err != nil {
		return fmt.Errorf("error setting password: %w", err)
	}

	params.Password = password
	err = s.repository.UpdatePassword(ctx, params)

	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	return nil
}

func (s *Service) UpdateActiveUser(ctx context.Context, params database.UpdateActiveParams) error {
	err := s.repository.UpdateActive(ctx, params)

	if err != nil {
		return fmt.Errorf("error updating status active user: %w", err)
	}

	return nil
}

func (s *Service) UpdateSuperUser(ctx context.Context, params database.UpdateSuperUserParams) error {
	err := s.repository.UpdateSuperUser(ctx, params)

	if err != nil {
		return fmt.Errorf("error updating status super user: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repository.DeleteUser(ctx, id)

	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

func (*Service) SetPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (*Service) MatchPassword(plainText string, password string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plainText, password)
	if err != nil {
		log.Fatal(err)
	}

	return match, nil
}
