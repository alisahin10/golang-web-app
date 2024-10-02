package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/utils"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/validator"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
	handler.log.Info("login endpoint called")
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req loginRequest
	if err := ctx.BodyParser(&req); err != nil {
		handler.log.Error("failed to parse body", zap.Error(err))
		return fiber.ErrBadRequest
	}
	if err := handler.validate.Struct(req); err != nil {
		handler.log.Error("failed to validate body", zap.Error(err))
		return fiber.ErrBadRequest
	}
	user, err := handler.repo.FindOneByEmail(req.Email)
	if err != nil {
		handler.log.Error("failed to find user by email", zap.Error(err))
		return fiber.ErrInternalServerError
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		handler.log.Error("invalid password", zap.Error(err))
		return fiber.ErrInternalServerError
	}
	accessToken, refreshToken, err := utils.GenerateTokens(user.Username, user.Email)
	if err != nil {
		handler.log.Error("failed to generate tokens", zap.Error(err))
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          utils.ToResponseUser(user),
	})

	return nil
}

func (handler *auth) logoutEndpoint(ctx *fiber.Ctx) error {
	token := ctx.Params("token")
	if token == "" {
		return fiber.ErrBadRequest
	}
	handler.log.Info("logout", zap.String("token", token))
	return ctx.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
	return nil
}
