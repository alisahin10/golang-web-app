package validator

import (
	"github.com/go-playground/validator/v10"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
	"regexp"
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
	v := validator.New()

	// Email format validation with regex
	v.RegisterValidation("email_format", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		return re.MatchString(fl.Field().String())
	})

	return &validatorImpl{
		v:    v,
		repo: repo,
	}
}

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
	if len(user.Password) < 8 {
		return false, "Password must be at least 8 characters long"
	}

	// Check if email is in the database or not
	existingUser, err := v.repo.FindOneByEmail(user.Email)
	if err == nil && existingUser != nil {
		return false, "Email already exists"
	}

	// Email format controller with regex validator
	if err := v.v.Var(user.Email, "required,email,email_format"); err != nil {
		return false, "Invalid email format"
	}

	// Age control for younger than 18
	if user.Age < 18 {
		return false, "User must be at least 18 years old"
	}

	return true, ""
}
