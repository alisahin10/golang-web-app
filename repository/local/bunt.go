package local

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/buntdb"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/utils"
)

// BuntImpl struct that holds the database instance
type BuntImpl struct {
	DB *buntdb.DB // BuntDB instance for database operations
}

// NewBuntRepository initializes a new BuntDB repository.
func NewBuntRepository(dbPath string) (Repository, error) {
	// Open a connection to the BuntDB database.
	db, err := buntdb.Open(dbPath)
	if err != nil {
		return nil, err // Return error if the database cannot be opened.
	}

	// Return a new instance of BuntImpl with the open database.
	return &BuntImpl{DB: db}, nil
}

// Create saves a new user to the database.
func (repo *BuntImpl) Create(user *model.User) error {
	return repo.DB.Update(func(tx *buntdb.Tx) error {
		// Convert user struct to JSON format for storage.
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err // Return error if JSON marshaling fails.
		}

		// Save the user data to the database.
		_, _, err = tx.Set(fmt.Sprintf("user:%s", user.ID), string(userJSON), nil)
		return err // Return any error encountered during save.
	})
}

// FindOneByID retrieves a user by their unique ID from the database.
func (repo *BuntImpl) FindOneByID(userID string) (*model.User, error) {
	var user model.User
	err := repo.DB.View(func(tx *buntdb.Tx) error {
		// Retrieve the user data from the database.
		val, err := tx.Get(fmt.Sprintf("user:%s", userID))
		if err != nil {
			return err // Return error if the user is not found.
		}
		// Unmarshal the JSON data into the user struct.
		return json.Unmarshal([]byte(val), &user)
	})
	if err != nil {
		return nil, err // Return error if fetching or unmarshalling fails.
	}
	return &user, nil // Return the retrieved user.
}

// FindAll retrieves all users from the database.
func (repo *BuntImpl) FindAll() ([]*model.User, error) {
	var users []*model.User // Slice to hold all users

	err := repo.DB.View(func(tx *buntdb.Tx) error {
		// Iterate through all user records in the database.
		err := tx.Ascend("", func(key, value string) bool {
			if len(key) > 5 && key[:5] == "user:" {
				var user model.User
				if err := json.Unmarshal([]byte(value), &user); err == nil {
					users = append(users, &user) // Append found users to the slice.
				}
			}
			return true // Continue iteration.
		})
		return err // Return any error encountered during iteration.
	})

	if err != nil {
		return nil, err // Return error if fetching fails.
	}
	return users, nil // Return the slice of users.
}

// UpdateOneByID updates user data for a given user ID.
func (repo *BuntImpl) UpdateOneByID(userID string, updateData *model.User) error {
	// Get current user from the database.
	user, err := repo.FindOneByID(userID)
	if err != nil {
		return err // Return error if user not found.
	}

	// Update the desired fields in the user struct.
	user.UpdateFields(updateData)

	// Hashing the password if it's changed.
	if updateData.Password != "" {
		hashedPassword, err := utils.HashPassword(updateData.Password)
		if err != nil {
			return fmt.Errorf("password hash error: %v", err) // Return error if password hashing fails.
		}
		user.Password = hashedPassword // Update user password.
	}

	// Marshal the updated user data into JSON format.
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("update user data JSON error: %v", err) // Return error if marshaling fails.
	}

	// Save the updated user data to the database.
	err = repo.DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(fmt.Sprintf("user:%s", userID), string(userJSON), nil)
		return err // Return any error encountered during save.
	})

	return err // Return any error from the update operation.
}

// DeleteOneByID removes a user from the database by their ID.
func (repo *BuntImpl) DeleteOneByID(userID string) error {
	// Delete user from the database by ID.
	err := repo.DB.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(fmt.Sprintf("user:%s", userID))
		if err != nil {
			return fmt.Errorf("user not found or error deleting user: %w", err) // Return error if user not found or delete fails.
		}
		return nil // Return nil if successful.
	})

	return err // Return any error from the delete operation.
}

// FindOneByEmail retrieves a user by their email address from the database.
func (repo *BuntImpl) FindOneByEmail(email string) (*model.User, error) {
	var user model.User

	// Search the user in db with email
	err := repo.DB.View(func(tx *buntdb.Tx) error {
		// Iterate all the users in the database.
		err := tx.Ascend("", func(key, value string) bool {
			if len(key) > 5 && key[:5] == "user:" {
				var u model.User
				if err := json.Unmarshal([]byte(value), &u); err == nil {
					if u.Email == email {
						user = u
						return false // Stop iteration when user found
					}
				}
			}
			return true // Continue iteration.
		})
		return err // Return any error encountered during iteration.
	})

	if err != nil {
		return nil, err // Return error if fetching fails.
	}

	// If user is not found, return nil
	if user.ID == "" {
		return nil, fmt.Errorf("user not found") // Return error if user ID is empty.
	}

	return &user, nil // Return the found user.
}

// SaveRefreshToken stores a user's refresh token in the database.
func (repo *BuntImpl) SaveRefreshToken(UserID string, refreshToken string) error {
	return repo.DB.Update(func(tx *buntdb.Tx) error {
		key := fmt.Sprintf("refresh_token:%s", UserID)
		_, _, err := tx.Set(key, refreshToken, nil)
		return err // Return any error encountered during save.
	})
}

// FindRefreshToken retrieves a user ID based on the provided refresh token.
func (repo *BuntImpl) FindRefreshToken(token string) (string, error) {
	var userID string
	err := repo.DB.View(func(tx *buntdb.Tx) error {
		// Iterate over keys that start with "refresh_token:"
		err := tx.Ascend("", func(key, value string) bool {
			if key[:14] == "refresh_token:" && value == token {
				// Extract the userID from the key
				userID = key[len("refresh_token:"):]
				return false // Stop iteration once the token is found
			}
			return true // Continue iteration.
		})
		return err // Return any error encountered during iteration.
	})
	if err != nil {
		return "", err // Return error if fetching fails.
	}
	if userID == "" {
		return "", fmt.Errorf("refresh token not found") // Return error if user ID is empty.
	}
	return userID, nil // Return the found user ID.
}

// DeleteRefreshToken removes a user's refresh token from the database.
func (repo *BuntImpl) DeleteRefreshToken(userID string) error {
	return repo.DB.Update(func(tx *buntdb.Tx) error {
		key := fmt.Sprintf("refresh_token:%s", userID)
		_, err := tx.Delete(key)
		return err // Return any error encountered during delete.
	})
}

// Close closes the database connection.
func (repo *BuntImpl) Close() error {
	return repo.DB.Close() // Return error if closing fails.
}
