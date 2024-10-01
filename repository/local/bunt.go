package local

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/buntdb"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/utils"
)

type BuntImpl struct {
	DB *buntdb.DB
}

func NewBuntRepository(dbPath string) (Repository, error) {
	// DB connection and error handling.
	db, err := buntdb.Open(dbPath)
	if err != nil {
		return nil, err
	}

	// Create buntImp repository instance
	return &BuntImpl{DB: db}, nil
}

func (repo *BuntImpl) Create(user *model.User) error {
	return repo.DB.Update(func(tx *buntdb.Tx) error {
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}

		// Save the user data.
		_, _, err = tx.Set(fmt.Sprintf("user:%s", user.ID), string(userJSON), nil)
		return err
	})
}

func (repo *BuntImpl) FindOneByID(userID string) (*model.User, error) {
	var user model.User
	err := repo.DB.View(func(tx *buntdb.Tx) error {
		// user:<userID> formatted data receive
		val, err := tx.Get(fmt.Sprintf("user:%s", userID))
		if err != nil {
			return err
		}
		// JSON conversion to received data
		return json.Unmarshal([]byte(val), &user)
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *BuntImpl) FindAll() ([]*model.User, error) {
	var users []*model.User

	err := repo.DB.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			if len(key) > 5 && key[:5] == "user:" { // "user:" prefix'ini kontrol ediyoruz
				var user model.User
				if err := json.Unmarshal([]byte(value), &user); err == nil {
					users = append(users, &user)
				}
			}
			return true
		})
		return err
	})

	if err != nil {
		return nil, err
	}
	return users, nil
}
func (repo *BuntImpl) UpdateOneByID(userID string, updateData *model.User) error {
	// Get current user from database.
	user, err := repo.FindOneByID(userID)
	if err != nil {
		return err
	}

	// Checking the desired update area.
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Lastname != "" {
		user.Lastname = updateData.Lastname
	}
	if updateData.Age != 0 {
		user.Age = updateData.Age
	}

	// Hashing the password if it's changed.
	if updateData.Password != "" {
		hashedPassword, err := utils.HashPassword(updateData.Password)
		if err != nil {
			return fmt.Errorf("password hash error: %v", err)
		}
		user.Password = hashedPassword
	}

	// Updating the user data to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("update user data JSON error: %v", err)
	}

	// Save it to the database
	err = repo.DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(fmt.Sprintf("user:%s", userID), string(userJSON), nil)
		return err
	})

	return err
}

/*
func (repo *BuntImpl) UpdateOneByID(userID, email, name, lastname string, age int) error {
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
	err = repo.DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(fmt.Sprintf("user:%s", userID), string(userJSON), nil)
		return err
	})
	return err
}
*/

func (repo *BuntImpl) DeleteOneByID(userID string) error {
	// Delete the user by ID
	err := repo.DB.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(fmt.Sprintf("user:%s", userID))
		return err
	})

	return err
}

func (repo *BuntImpl) FindOneByEmail(email string) (*model.User, error) {
	var user model.User

	// Search the user in db with email
	err := repo.DB.View(func(tx *buntdb.Tx) error {
		// Iterate all the users
		err := tx.Ascend("", func(key, value string) bool {
			if len(key) > 5 && key[:5] == "user:" {
				var u model.User
				if err := json.Unmarshal([]byte(value), &u); err == nil {
					if u.Email == email {
						user = u
						return false // Stop the iteration when user found
					}
				}
			}
			return true
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	// If user is not found then return nil
	if user.ID == "" {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}
