package validator

import (
	"github.com/go-playground/validator/v10"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"regexp"
)

// Validate interface for the validator.
type Validate interface {
	Struct(s interface{}) error
	ValidateUser(user *model.User) (bool, string)
	ValidateEmailFormat(email string) bool
}

type validatorImpl struct {
	v *validator.Validate
}

// NewValidator creates a new instance of validator.
func NewValidator() Validate {
	v := validator.New()

	// Email format validation with regex
	v.RegisterValidation("email_format", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		return re.MatchString(fl.Field().String())
	})

	return &validatorImpl{v: v}
}

// Struct performs validation on a struct.
func (v *validatorImpl) Struct(s interface{}) error {
	return v.v.Struct(s)
}

// ValidateUser validates the user struct.
func (v *validatorImpl) ValidateUser(user *model.User) (bool, string) {
	// Validation logic...
	return true, ""
}

// ValidateEmailFormat checks the format of the email.
func (v *validatorImpl) ValidateEmailFormat(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
