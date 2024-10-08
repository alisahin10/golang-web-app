package middleware

import (
	"github.com/gofiber/fiber/v2"
)

type AppError struct{}

// NewBadRequest returns a 400 Bad Request error with a custom message
func (e *AppError) NewBadRequest(message string) *fiber.Error {
	return &fiber.Error{
		Code:    fiber.StatusBadRequest,
		Message: message,
	}
}

// NewUnauthorized returns a 401 Unauthorized error with a custom message
func (e *AppError) NewUnauthorized(message string) *fiber.Error {
	return &fiber.Error{
		Code:    fiber.StatusUnauthorized,
		Message: message,
	}
}

// NewInternalServerError returns a 500 Internal Server Error with a custom message
func (e *AppError) NewInternalServerError(message string) *fiber.Error {
	return &fiber.Error{
		Code:    fiber.StatusInternalServerError,
		Message: message,
	}
}

// NewNotFound returns a 404 Not Found error with a custom message
func (e *AppError) NewNotFound(message string) *fiber.Error {
	return &fiber.Error{
		Code:    fiber.StatusNotFound,
		Message: message,
	}
}
