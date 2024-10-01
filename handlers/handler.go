package handlers

import "github.com/gofiber/fiber/v2"

type Handler interface {
	AssignEndpoints(prefix string, router fiber.Router)
}

type UserHandler interface {
	AssignUserEndpoints(prefix string, router fiber.Router)
}
