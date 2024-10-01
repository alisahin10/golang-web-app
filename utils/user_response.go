package utils

import "gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"

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

// ToResponseUsers converts a list of User models to a list of UserResponse models
func ToResponseUsers(users []model.User) []model.UserResponse {
	var responseUsers []model.UserResponse
	for _, user := range users {
		responseUsers = append(responseUsers, ToResponseUser(&user))
	}
	return responseUsers
}
