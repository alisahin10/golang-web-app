package services

import (
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
)

// UserService defines the interface for user-related operations.
type UserService interface {
	IsEmailTaken(email string) (bool, error)
	FindByEmail(email string) (*model.User, error)
}

type userServiceImpl struct {
	repo local.Repository
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo local.Repository) UserService {
	return &userServiceImpl{repo: repo}
}

// IsEmailTaken checks if an email is already taken.
func (s *userServiceImpl) IsEmailTaken(email string) (bool, error) {
	user, err := s.repo.FindOneByEmail(email)
	if err != nil {
		// If user is not found, it means the email is not taken.
		if err.Error() == "user not found" {
			return false, nil
		}
		return false, err
	}
	// If user is found, email is taken.
	return user != nil, nil
}

// Implementation of FindByEmail
func (s *userServiceImpl) FindByEmail(email string) (*model.User, error) {
	user, err := s.repo.FindOneByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
