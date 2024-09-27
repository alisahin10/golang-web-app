package local

import "gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"

type buntImpl struct {
}

func NewBuntRepository(dbPAth string) (Repository, error) {
	// Create bunt db
	return &buntImpl{}, nil
}

func (repo *buntImpl) Create(user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (repo *buntImpl) FindOneByID(userID string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *buntImpl) FindAll() ([]*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *buntImpl) UpdateOneByID(userID, email, name, lastname string, age int) error {
	//TODO implement me
	panic("implement me")
}

func (repo *buntImpl) DeleteOneByID(userID string) error {
	//TODO implement me
	panic("implement me")
}
