package model

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Age      int    `json:"age"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Age      int    `json:"age"`
}

// Display the response in order for Create function.
type CreateUserResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Lastname     string `json:"lastname"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Update user function encapsulation.
// UpdateFields updates non-zero fields of the current user with the update data
func (u *User) UpdateFields(updateData *User) {
	if updateData.Username != "" {
		u.Username = updateData.Username
	}
	if updateData.Email != "" {
		u.Email = updateData.Email
	}
	if updateData.Name != "" {
		u.Name = updateData.Name
	}
	if updateData.Lastname != "" {
		u.Lastname = updateData.Lastname
	}
	if updateData.Age != 0 {
		u.Age = updateData.Age
	}
	if updateData.Password != "" {
		u.Password = updateData.Password // Hashing will still be done in the caller function
	}
}
