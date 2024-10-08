package services

import (
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
)

// UserService defines the interface for user-related operations.
type UserService interface {
	IsEmailTaken(email string) (bool, error)       // Check if an email is already taken.
	FindByEmail(email string) (*model.User, error) // Retrieve a user by their email address.
}

type userServiceImpl struct {
	repo local.Repository // Reference to the repository for database operations.
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo local.Repository) UserService {
	return &userServiceImpl{repo: repo} // Initialize userServiceImpl with the provided repository.
}

// IsEmailTaken checks if an email is already taken.
func (s *userServiceImpl) IsEmailTaken(email string) (bool, error) {
	// Attempt to find a user by the provided email.
	user, err := s.repo.FindOneByEmail(email)
	if err != nil {
		// If user is not found, it means the email is not taken.
		if err.Error() == "user not found" {
			return false, nil // Email is not taken, return false.
		}
		// If another error occurred while checking, return it.
		return false, err
	}
	// If user is found, email is taken.
	return user != nil, nil // Return true if the user exists (email is taken).
}

// FindByEmail retrieves a user by their email address.
func (s *userServiceImpl) FindByEmail(email string) (*model.User, error) {
	// Attempt to find a user by the provided email.
	user, err := s.repo.FindOneByEmail(email)
	if err != nil {
		return nil, err // Return error if user retrieval fails.
	}
	return user, nil // Return the found user.
}
