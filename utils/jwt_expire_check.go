package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// IsExpired checks whether a given JWT token has expired.
func IsExpired(tokenStr string, jwtSecret []byte) bool {
	// Parse the token using the JWT secret
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil // Use the secret passed as a parameter
	})

	// If there's an error or token is invalid, assume it's expired
	if err != nil || !token.Valid {
		return true
	}

	// Extract the exp claim (expiration time)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp := int64(claims["exp"].(float64)) // Convert exp to int64
		return time.Now().Unix() > exp        // Check if current time is past the expiration time
	}

	return true // If claims are missing or invalid, assume token is expired
}
