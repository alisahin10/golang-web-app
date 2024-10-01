package local

import (
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
)

type Repository interface {
	// Typed as instance
	Create(user *model.User) error
	FindOneByID(userID string) (*model.User, error)
	FindAll() ([]*model.User, error)
	UpdateOneByID(userID, email, name, lastname string, age int) error
	DeleteOneByID(userID string) error
	FindOneByEmail(email string) (*model.User, error)
}
