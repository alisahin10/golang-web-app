package validator

import (
	"github.com/go-playground/validator/v10"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
)

type Validate interface {
	Struct(s interface{}) error
	ValidateUser(user *model.User) (bool, string)
}

type validatorImpl struct {
	v    *validator.Validate
	repo local.Repository
}

func NewValidator(repo local.Repository) Validate {
	return &validatorImpl{
		v:    validator.New(),
		repo: repo,
	}
}

/*
func NewValidator() Validate {
	return &validatorImpl{v: validator.New()}
}

*/

func (v *validatorImpl) Struct(s interface{}) error {
	return v.v.Struct(s)
}

// ValidateUser validates the necessary fields of the user object.
func (v *validatorImpl) ValidateUser(user *model.User) (bool, string) {
	if user.Name == "" {
		return false, "Name is required"
	}
	if user.Lastname == "" {
		return false, "Lastname is required"
	}
	if user.Username == "" {
		return false, "Username is required"
	}
	if user.Email == "" {
		return false, "Email is required"
	}
	if user.Password == "" {
		return false, "Password is required"
	}

	// Email'in veritabanında olup olmadığını kontrol ediyoruz
	existingUser, err := v.repo.FindOneByEmail(user.Email)
	if err == nil && existingUser != nil {
		return false, "Email already exists"
	}

	return true, ""
}
