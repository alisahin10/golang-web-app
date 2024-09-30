package local

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/buntdb"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"os"
)

type buntImpl struct {
	db *buntdb.DB
}

func NewBuntRepository(dbPath string) (Repository, error) {
	db, err := buntdb.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	localDbPath := os.Getenv("LOCAL_DB_PATH")
	if _, err = os.Stat(localDbPath); os.IsNotExist(err) {
		err := os.MkdirAll(localDbPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return &buntImpl{}, nil
}

func (repo *buntImpl) Create(user *model.User) error {
	return repo.db.Update(func(tx *buntdb.Tx) error {
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}
		err = repo.db.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(fmt.Sprintf("user:%s", user.ID), string(userJSON), nil)
			return err
		})

		if err != nil {
			return err
		}
		return nil
	})

}

func (repo *buntImpl) FindOneByID(userID string) (*model.User, error) {
	var user model.User

	// Read the user data from BuntDB
	err := repo.db.View(func(tx *buntdb.Tx) error {
		userJSON, err := tx.Get(fmt.Sprintf("user:%s", userID))
		if err != nil {
			return err
		}

		// Unmarshal JSON string back into User struct
		err = json.Unmarshal([]byte(userJSON), &user)
		return err
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *buntImpl) FindAll() ([]*model.User, error) {
	var users []*model.User

	// Iterate through all users in the database
	err := repo.db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			if len(key) > 5 && key[:5] == "user:" { // Check for keys with "user:" prefix
				var user model.User
				if err := json.Unmarshal([]byte(value), &user); err == nil {
					users = append(users, &user)
				}
			}
			return true // continue iteration
		})
		return err
	})

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *buntImpl) UpdateOneByID(userID, email, name, lastname string, age int) error {
	// Fetch the existing user
	user, err := repo.FindOneByID(userID)
	if err != nil {
		return err
	}
	user.Email = email
	user.Name = name
	user.Lastname = lastname
	user.Age = age
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = repo.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(fmt.Sprintf("user:%s", userID), string(userJSON), nil)
		return err
	})
	return err
}

func (repo *buntImpl) DeleteOneByID(userID string) error {
	// Delete the user by ID
	err := repo.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(fmt.Sprintf("user:%s", userID))
		return err
	})

	return err
}
