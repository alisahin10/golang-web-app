package model

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	Token string `json:"token" `
}
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	Identifier   string `json:"identifier"`
}
