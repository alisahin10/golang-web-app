package utils

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // JWT_SECRET environment variable

// GenerateTokens creates both access and refresh tokens for a user
func GenerateTokens(username, email string) (string, string, error) {
	// Access token for 10 minutes.
	accessTokenClaims := jwt.MapClaims{
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh token for 7 days.
	refreshTokenClaims := jwt.MapClaims{
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
