package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var errors *AppError

// JWTAuthMiddleware verifies JWT token and authorizes users for protected routes
func JWTAuthMiddleware(jwtSecret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// JWT validation logic using jwtSecret
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return errors.NewUnauthorized("Unauthorized, no token provided")
		}

		// Remove "Bearer " prefix
		tokenString = tokenString[len("Bearer "):]

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		// Validate token
		if err != nil || !token.Valid {
			return errors.NewUnauthorized("Unauthorized, invalid token")
		}

		// Extract claims and proceed
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return errors.NewUnauthorized("Unauthorized, invalid token")
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
