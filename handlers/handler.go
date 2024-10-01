package handlers

import "github.com/gofiber/fiber/v2"

type Handler interface {
	AssignEndpoints(prefix string, router fiber.Router)
	AssignUserEndpoints(prefix string, router fiber.Router)
}
