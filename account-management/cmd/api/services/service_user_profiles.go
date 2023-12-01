package services

import "account-management/internal/database"

type UserProfileService struct {
	repository *database.Queries
}

func NewUserProfileService(rep *database.Queries) *UserProfileService {
	return &UserProfileService{
		repository: rep,
	}
}
