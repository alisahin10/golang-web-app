package utils

import (
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
)

// ToResponseUser converts a User model to a UserResponse model
func ToResponseUser(user *model.User) model.UserResponse {
	return model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Lastname: user.Lastname,
		Email:    user.Email,
		Age:      user.Age,
	}
}

// ToCreateUserResponse generates a CreateUserResponse from a User model and tokens.
func ToCreateUserResponse(user *model.User, accessToken string, refreshToken string) model.CreateUserResponse {
	return model.CreateUserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Name:         user.Name,
		Lastname:     user.Lastname,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
