package local

import (
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"os"
	"testing"
)

func TestCreateUser(t *testing.T) {
	os.Setenv("LOCAL_DB_PATH", "./test.db")
	defer os.Remove("./test.db")

	repo, err := NewBuntRepository(os.Getenv("LOCAL_DB_PATH"))
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	testCases := []struct {
		name    string
		user    model.User
		wantErr bool
	}{
		{
			name:    "Valid User",
			user:    model.User{ID: "1", Username: "testuser", Email: "test@example.com", Password: "password123", Name: "Test", Lastname: "User", Age: 25},
			wantErr: false,
		},
		{
			name:    "Empty Fields",
			user:    model.User{ID: "2", Username: "", Email: "", Password: "", Name: "", Lastname: "", Age: 0},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Create(&tc.user)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Create() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestFindOneByID(t *testing.T) {
	os.Setenv("LOCAL_DB_PATH", "./test_find.db")
	defer os.Remove("./test_find.db")

	repo, _ := NewBuntRepository(os.Getenv("LOCAL_DB_PATH"))

	// Set up a known user
	_ = repo.Create(&model.User{ID: "123", Email: "test@example.com"})

	testCases := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "User Exists",
			userID:  "123",
			wantErr: false,
		},
		{
			name:    "User Does Not Exist",
			userID:  "999",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.FindOneByID(tc.userID)
			if (err != nil) != tc.wantErr {
				t.Fatalf("FindOneByID() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}
func TestUpdateUser(t *testing.T) {
	os.Setenv("LOCAL_DB_PATH", "./test_update.db")
	defer os.Remove("./test_update.db")

	repo, _ := NewBuntRepository(os.Getenv("LOCAL_DB_PATH"))

	// Setup: Create user
	_ = repo.Create(&model.User{ID: "123", Email: "old@example.com", Password: "oldpass"})

	testCases := []struct {
		name       string
		userID     string
		updateData model.User
		wantErr    bool
	}{
		{
			name:       "Update Email and Password",
			userID:     "123",
			updateData: model.User{Email: "new@example.com", Password: "newpass"},
			wantErr:    false,
		},
		{
			name:       "Update Non-Existent User",
			userID:     "999",
			updateData: model.User{Email: "new@example.com"},
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.UpdateOneByID(tc.userID, &tc.updateData)
			if (err != nil) != tc.wantErr {
				t.Fatalf("UpdateOneByID() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}
func TestDeleteUser(t *testing.T) {
	os.Setenv("LOCAL_DB_PATH", "./test_delete.db")
	defer os.Remove("./test_delete.db")

	repo, _ := NewBuntRepository(os.Getenv("LOCAL_DB_PATH"))

	// Setup: Create user
	_ = repo.Create(&model.User{ID: "123"})

	testCases := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "Delete Existing User",
			userID:  "123",
			wantErr: false,
		},
		{
			name:    "Delete Non-Existent User",
			userID:  "999",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.DeleteOneByID(tc.userID)
			if (err != nil) != tc.wantErr {
				t.Fatalf("DeleteOneByID() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}
func TestFindOneByEmail(t *testing.T) {
	os.Setenv("LOCAL_DB_PATH", "./test_find_email.db")
	defer os.Remove("./test_find_email.db")

	repo, _ := NewBuntRepository(os.Getenv("LOCAL_DB_PATH"))

	// Setup: Create user
	_ = repo.Create(&model.User{ID: "123", Email: "test@example.com"})

	testCases := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Find Existing User by Email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Find Non-Existent User by Email",
			email:   "nonexistent@example.com",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.FindOneByEmail(tc.email)
			if (err != nil) != tc.wantErr {
				t.Fatalf("FindOneByEmail() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestFindAll(t *testing.T) {
	// Set up environment for test
	os.Setenv("LOCAL_DB_PATH", "./test_find_all.db")
	defer os.Remove("./test_find_all.db")

	// Create a new repository
	repo, err := NewBuntRepository(os.Getenv("LOCAL_DB_PATH"))
	if err != nil {
		t.Fatalf("Error creating repository: %v", err)
	}

	// Add multiple users to the repository
	users := []*model.User{
		{ID: "1", Username: "user1", Email: "user1@example.com"},
		{ID: "2", Username: "user2", Email: "user2@example.com"},
		{ID: "3", Username: "user3", Email: "user3@example.com"},
	}

	// Insert users into the database
	for _, user := range users {
		err := repo.Create(user)
		if err != nil {
			t.Fatalf("Error creating user %s: %v", user.Username, err)
		}
	}

	// Call the FindAll function
	foundUsers, err := repo.FindAll()
	if err != nil {
		t.Fatalf("Error finding all users: %v", err)
	}

	// Verify that all users are retrieved
	if len(foundUsers) != len(users) {
		t.Fatalf("Expected %d users, found %d", len(users), len(foundUsers))
	}

	// Map the found users by ID for easier comparison
	foundUserMap := make(map[string]*model.User)
	for _, user := range foundUsers {
		foundUserMap[user.ID] = user
	}

	// Check that all created users are found
	for _, user := range users {
		foundUser, exists := foundUserMap[user.ID]
		if !exists {
			t.Fatalf("User with ID %s not found", user.ID)
		}
		if foundUser.Username != user.Username || foundUser.Email != user.Email {
			t.Fatalf("Expected user %v, found %v", user, foundUser)
		}
	}
}

func TestSaveRefreshToken(t *testing.T) {
	os.Setenv("LOCAL_DB_PATH", "./test_refresh_token.db")
	defer os.Remove("./test_refresh_token.db")

	repo, err := NewBuntRepository(os.Getenv("LOCAL_DB_PATH"))
	if err != nil {
		t.Fatalf("Error creating repository: %v", err)
	}

	// Save a refresh token for a user
	err = repo.SaveRefreshToken("user123", "token123")
	if err != nil {
		t.Fatalf("Error saving refresh token: %v", err)
	}

	// Verify that the token was saved correctly
	tokenOwner, err := repo.FindRefreshToken("token123")
	if err != nil {
		t.Fatalf("Error finding refresh token: %v", err)
	}
	if tokenOwner != "user123" {
		t.Fatalf("Expected userID 'user123', got '%s'", tokenOwner)
	}
}

func TestFindRefreshToken(t *testing.T) {
	os.Setenv("LOCAL_DB_PATH", "./test_find_refresh_token.db")
	defer os.Remove("./test_find_refresh_token.db")

	repo, err := NewBuntRepository(os.Getenv("LOCAL_DB_PATH"))
	if err != nil {
		t.Fatalf("Error creating repository: %v", err)
	}

	// Save a refresh token
	_ = repo.SaveRefreshToken("user123", "token123")

	testCases := []struct {
		name       string
		token      string
		wantUserID string
		wantErr    bool
	}{
		{
			name:       "Token Exists",
			token:      "token123",
			wantUserID: "user123",
			wantErr:    false,
		},
		{
			name:       "Token Does Not Exist",
			token:      "token999",
			wantUserID: "",
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := repo.FindRefreshToken(tc.token)
			if (err != nil) != tc.wantErr {
				t.Fatalf("FindRefreshToken() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if userID != tc.wantUserID {
				t.Fatalf("FindRefreshToken() userID = %v, want %v", userID, tc.wantUserID)
			}
		})
	}
}

func TestDeleteRefreshToken(t *testing.T) {
	os.Setenv("LOCAL_DB_PATH", "./test_delete_refresh_token.db")
	defer os.Remove("./test_delete_refresh_token.db")

	repo, err := NewBuntRepository(os.Getenv("LOCAL_DB_PATH"))
	if err != nil {
		t.Fatalf("Error creating repository: %v", err)
	}

	// Save a refresh token
	_ = repo.SaveRefreshToken("user123", "token123")

	// Delete the refresh token
	err = repo.DeleteRefreshToken("user123")
	if err != nil {
		t.Fatalf("Error deleting refresh token: %v", err)
	}

	// Verify that the token was deleted
	_, err = repo.FindRefreshToken("token123")
	if err == nil {
		t.Fatalf("Expected error finding deleted refresh token, got nil")
	}
}
