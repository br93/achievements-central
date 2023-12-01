package main

import (
	"account-management/cmd/api/services"
)

type Services struct {
	User        *services.UserService
	UserProfile *services.UserProfileService
}
