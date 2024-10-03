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

func NewUser(log *zap.Logger, repo local.Repository, validate validator.Validate) Handler {
	return &user{
		log:      log,
		repo:     repo,
		validate: validate,
	}
}

func (handler *user) AssignEndpoints(prefix string, router fiber.Router) {
	r := router.Group(prefix)

	r.Post("create", handler.createEndpoint)
	r.Get("/search", handler.findByEmailEndpoint)
	r.Get(":id", handler.getEndpoint)
	r.Get("/", handler.getAllEndpoint)
	r.Patch("update/:id", handler.updateEndpoint)
	r.Delete("/:id", handler.deleteEndpoint)

}

func (handler *user) createEndpoint(c *fiber.Ctx) error {
	user := new(model.User)

	// Parse JSON body into user model
	if err := c.BodyParser(user); err != nil {
		handler.log.Error("Error parsing body", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate the user data using the ValidateUser method from the validator
	isValid, validationErr := handler.validate.ValidateUser(user)
	if !isValid {
		handler.log.Error("Validation error", zap.String("error", validationErr))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr})
	}

	// Hash the user's password.
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		handler.log.Error("Failed to hash password", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}
	user.Password = hashedPassword

	// Assign a UUID to the user using the UUID utility function.
	user.ID = utils.GenerateUUID()

	// Use the repository to create a new user
	if err := handler.repo.Create(user); err != nil {
		handler.log.Error("Error creating user", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
	}

	// Create JWT token.
	accessToken, refreshToken, err := utils.GenerateTokens(user.Username, user.Email)
	if err != nil {
		handler.log.Error("Failed to generate tokens", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate tokens")
	}

	// Save the refresh token in the database
	if err := handler.repo.SaveRefreshToken(user.ID, refreshToken); err != nil {
		handler.log.Error("Failed to save refresh token", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save refresh token")
	}

	// Use the utility function to generate the response
	response := utils.ToCreateUserResponse(user, accessToken, refreshToken)

	handler.log.Info("User created successfully", zap.String("userID", user.ID))
	return c.Status(http.StatusCreated).JSON(response)

}

func (handler *user) getEndpoint(c *fiber.Ctx) error {
	// Get user information via ID.
	userID := c.Params("id")
	handler.log.Info("UserID param:", zap.String("userID", userID))

	// Get user data from database.
	user, err := handler.repo.FindOneByID(userID)
	if err != nil {
		handler.log.Error("User not found in database", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if user fields are populated
	if user.Name == "" || user.Email == "" {
		handler.log.Error("User found but fields are empty", zap.String("userID", userID))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User data is incomplete"})
	}

	userResponse := utils.ToResponseUser(user)

	handler.log.Info("User found:", zap.String("username", user.Username), zap.String("email", user.Email))

	return c.Status(fiber.StatusOK).JSON(userResponse)
}

func (handler *user) getAllEndpoint(c *fiber.Ctx) error {
	// Take users from database
	users, err := handler.repo.FindAll()
	if err != nil {
		handler.log.Error("Error fetching users from database", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch users"})
	}

	// If there is no user then empty slice returned
	if len(users) == 0 {
		handler.log.Info("No users found in the database")
		return c.Status(fiber.StatusOK).JSON([]model.UserResponse{})
	}

	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = utils.ToResponseUser(user)
	}

	handler.log.Info("All users fetched successfully")
	return c.Status(fiber.StatusOK).JSON(userResponses)
}

func (handler *user) updateEndpoint(c *fiber.Ctx) error {
	// Receiving user ID.
	userID := c.Params("id")
	handler.log.Info("UserID param:", zap.String("userID", userID))

	// Receiving update data from body
	updateData := new(model.User)
	if err := c.BodyParser(updateData); err != nil {
		handler.log.Error("Error parsing update data", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Finding the user in database
	err := handler.repo.UpdateOneByID(userID, updateData)
	if err != nil {
		handler.log.Error("Error updating user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
	}

	handler.log.Info("User updated successfully", zap.String("userID", userID))

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user_id": userID,
	})
}

func (handler *user) deleteEndpoint(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Check if user available before delete
	_, err := handler.repo.FindOneByID(userID)
	if err != nil {
		handler.log.Error("User not found", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Delete the user
	if err := handler.repo.DeleteOneByID(userID); err != nil {
		handler.log.Error("Error deleting user", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Error deleting user")
	}

	handler.log.Info("User deleted successfully", zap.String("userID", userID))

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
		"user_id": userID,
	})
}

func (handler *user) findByEmailEndpoint(c *fiber.Ctx) error {
	// Take query parameters from email
	email := c.Query("email")

	// Email parameter control
	if email == "" {
		handler.log.Error("Email query parameter is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email query parameter is required"})
	}

	// Verifying the mail receive with logging
	handler.log.Info("Searching for user with email:", zap.String("email", email))

	// // Search the user in the database via email
	user, err := handler.repo.FindOneByEmail(email)
	if err != nil {
		handler.log.Error("User not found by email", zap.String("email", email), zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// User found response message with JSON
	userResponse := utils.ToResponseUser(user)

	handler.log.Info("User found by email", zap.String("email", email))
	return c.Status(fiber.StatusOK).JSON(userResponse)
}
