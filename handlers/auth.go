package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/utils"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/validator"
	"go.uber.org/zap"
	"net/http"
)

type auth struct {
	log      *zap.Logger
	repo     local.Repository
	validate validator.Validate
}

func NewAuth(log *zap.Logger, repo local.Repository, validate validator.Validate) Handler {
	return &auth{
		log:      log,
		repo:     repo,
		validate: validate,
	}
}

func (handler *auth) AssignEndpoints(prefix string, router fiber.Router) {
	r := router.Group(prefix)

	r.Post("register", handler.registerEndpoint)
	r.Post("login", handler.loginEndpoint)
	r.Post("logout", handler.logoutEndpoint)
}

func (h *auth) registerEndpoint(c *fiber.Ctx) error {
	user := new(model.User)

	// Parse JSON body into user model
	if err := c.BodyParser(user); err != nil {
		h.log.Error("Error parsing body", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate the user data using the ValidateUser method from the validator
	isValid, validationErr := h.validate.ValidateUser(user)
	if !isValid {
		h.log.Error("Validation error", zap.String("error", validationErr))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr})
	}

	// Hash the user's password.
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		h.log.Error("Failed to hash password", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}
	user.Password = hashedPassword

	// Assign a UUID to the user using the UUID utility function.
	user.ID = utils.GenerateUUID()

	// Use the repository to create a new user
	if err := h.repo.Create(user); err != nil {
		h.log.Error("Error creating user", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
	}

	// JWT Token'ları oluştur
	accessToken, refreshToken, err := utils.GenerateTokens(user.Username, user.Email)
	if err != nil {
		h.log.Error("Failed to generate tokens", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate tokens")
	}

	//Get user in response and return tokens.
	response := fiber.Map{
		"username":      user.Username,
		"email":         user.Email,
		"name":          user.Name,
		"lastname":      user.Lastname,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	h.log.Info("User created successfully", zap.String("userID", user.ID))
	return c.Status(http.StatusCreated).JSON(response)

}

func (handler *auth) loginEndpoint(ctx *fiber.Ctx) error {
	// TODO implement me
	return nil
}

func (handler *auth) logoutEndpoint(ctx *fiber.Ctx) error {
	// TODO implement me
	return nil
}
