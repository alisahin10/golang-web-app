package validator

import "github.com/go-playground/validator/v10"

type Validate interface {
	Struct(s interface{}) error
}

type validatorImpl struct {
	v *validator.Validate
}

func NewValidator() Validate {
	return &validatorImpl{v: validator.New()}
}

func (v *validatorImpl) Struct(s interface{}) error {
	return v.v.Struct(s)
}
