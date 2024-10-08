package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// JWTAuthMiddleware verifies JWT token and authorizes users for protected routes
func JWTAuthMiddleware(c *fiber.Ctx) error {
	// Get the Authorization header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, no token provided",
		})
	}

	// Remove the "Bearer " prefix from the token
	tokenString = tokenString[len("Bearer "):]

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// If parsing or validation fails, return Unauthorized
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, invalid token",
		})
	}

	// Extract user ID and role from the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, invalid token claims",
		})
	}

	// Set user_id and role in the request context for later use
	c.Locals("user_id", claims["user_id"])
	c.Locals("role", claims["role"])

	// Continue to the next handler
	return c.Next()
}
