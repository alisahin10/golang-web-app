package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/validator"
	"go.uber.org/zap"
)

type auth struct {
	log      *zap.Logger
	repo     local.Repository
	validate validator.Validate
}

func NewAuth(log *zap.Logger, repo local.Repository, validate validator.Validate) *auth {
	return &auth{
		log:      log,
		repo:     repo,
		validate: validate,
	}
}

func (handler *auth) AssignEndpoints(prefix string, router fiber.Router) {
	r := router.Group(prefix)

	r.Post("login", handler.loginEndpoint)
	r.Post("logout", handler.logoutEndpoint)
}

func (handler *auth) loginEndpoint(ctx *fiber.Ctx) error {
	// TODO implement me
	return nil
}

func (handler *auth) logoutEndpoint(ctx *fiber.Ctx) error {
	// TODO implement me
	return nil
}
