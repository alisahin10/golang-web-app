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

type user struct {
	log      *zap.Logger
	repo     local.Repository
	validate validator.Validate
}

func NewUser(log *zap.Logger, repo local.Repository, validate validator.Validate) *user {
	return &user{
		log:      log,
		repo:     repo,
		validate: validate,
	}
}

func (handler *user) AssignUserEndpoints(prefix string, router fiber.Router) {
	r := router.Group(prefix)

	r.Post("create", handler.createEndpoint)
	r.Get(":id", handler.getEndpoint)
	r.Get("/", handler.getAllEndpoint)
}

func (h *user) createEndpoint(c *fiber.Ctx) error {
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

	// Create JWT token.
	accessToken, refreshToken, err := utils.GenerateTokens(user.Username, user.Email)
	if err != nil {
		h.log.Error("Failed to generate tokens", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate tokens")
	}

	//Get user in response and return tokens.
	response := fiber.Map{
		"id":            user.ID,
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

func (h *user) getEndpoint(c *fiber.Ctx) error {
	// Get user information via ID.
	userID := c.Params("id")
	h.log.Info("UserID param:", zap.String("userID", userID))

	// Get user data from database.
	user, err := h.repo.FindOneByID(userID)
	if err != nil {
		h.log.Error("User not found in database", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if user fields are populated
	if user.Name == "" || user.Email == "" {
		h.log.Error("User found but fields are empty", zap.String("userID", userID))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User data is incomplete"})
	}

	userResponse := utils.ToResponseUser(user)

	h.log.Info("User found:", zap.String("username", user.Username), zap.String("email", user.Email))

	return c.Status(fiber.StatusOK).JSON(userResponse)
}

func (h *user) getAllEndpoint(c *fiber.Ctx) error {
	// Take users from database
	users, err := h.repo.FindAll()
	if err != nil {
		h.log.Error("Error fetching users from database", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch users"})
	}

	// If there is no user then empty slice returned
	if len(users) == 0 {
		h.log.Info("No users found in the database")
		return c.Status(fiber.StatusOK).JSON([]model.UserResponse{})
	}

	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = utils.ToResponseUser(user)
	}

	h.log.Info("All users fetched successfully")
	return c.Status(fiber.StatusOK).JSON(userResponses)
}
