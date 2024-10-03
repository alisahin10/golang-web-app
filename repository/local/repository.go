package local

import (
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
)

type Repository interface {
	// Typed as instance
	Create(user *model.User) error
	FindOneByID(userID string) (*model.User, error)
	FindAll() ([]*model.User, error)
	UpdateOneByID(userID string, updateData *model.User) error
	DeleteOneByID(userID string) error
	FindOneByEmail(email string) (*model.User, error)
	SaveRefreshToken(UserID string, refreshToken string) error
	FindRefreshToken(UserID string) (string, error)
	DeleteRefreshToken(userID string) error
}
