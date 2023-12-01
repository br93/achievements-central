package services

import "account-management/internal/database"

type UserService struct {
	repository *database.Queries
}

func NewUserService(rep *database.Queries) *UserService {
	return &UserService{
		repository: rep,
	}
}
