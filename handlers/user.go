package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/middleware"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/services"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/utils/id"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/utils/jwt"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/utils/password"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/validator"
	"go.uber.org/zap"
	"net/http"
)

type user struct {
	log         *zap.Logger
	repo        local.Repository
	validate    validator.Validate
	userService services.UserService
	config      *AppConfig
}

// NewUser initializes a new user handler with dependencies.
func NewUser(log *zap.Logger, repo local.Repository, validate validator.Validate, config *AppConfig, userService services.UserService) Handler {
	return &user{
		log:         log,
		repo:        repo,
		validate:    validate,
		config:      config,
		userService: userService,
	}
}

// AssignEndpoints sets up the routes for the user-related operations.
func (handler *user) AssignEndpoints(prefix string, router fiber.Router) {
	r := router.Group(prefix)

	// Route for user creation, no JWT middleware here
	r.Post("create", handler.createEndpoint) // POST /user/create: Creates a new user and returns JWT tokens.

	// Routes that don't require authentication
	r.Get("/search", handler.findByEmailEndpoint) // GET /user/search: Searches for a user by email.
	r.Get(":id", handler.getEndpoint)             // GET /user/:id: Retrieves user information by ID.
	r.Get("/", handler.getAllEndpoint)            // GET /user: Retrieves a list of all users.

	// Routes that require JWT authentication
	//protectedRoutes := r.Group("/", middleware.JWTAuthMiddleware)

	// Routes that require JWT authentication
	protectedRoutes := r.Group("/", middleware.JWTAuthMiddleware(handler.config.JWTSecret))

	// These routes require the user to be authenticated (JWT)
	protectedRoutes.Patch("update/:id", handler.updateEndpoint) // PATCH /user/update/:id: Updates user information.
	protectedRoutes.Delete("/:id", handler.deleteEndpoint)      // DELETE /user/:id: Deletes a user by ID.
}

// createEndpoint handles user creation and returns JWT tokens upon success.
func (handler *user) createEndpoint(c *fiber.Ctx) error {
	user := new(model.User)

	// Parse JSON body into user model
	if err := c.BodyParser(user); err != nil {
		handler.log.Error("Error parsing body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate the user data using the ValidateUser method from the validator
	isValid, validationErr := handler.validate.ValidateUser(user)
	if !isValid {
		handler.log.Error("Validation error", zap.String("error", validationErr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationErr})
	}

	// Check if email is already taken using the UserService
	emailTaken, err := handler.userService.IsEmailTaken(user.Email)
	if err != nil {
		handler.log.Error("Error checking email", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not check email"})
	}
	if emailTaken {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is already taken"})
	}

	// Hash the user's password.
	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		handler.log.Error("Failed to hash password", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}
	user.Password = hashedPassword

	// Assign a UUID to the user using the UUID utility function.
	user.ID = id.GenerateUUID()

	// Assign default role to the new user
	user.Role = "user" // Default role assigned to new users

	// Use the repository to create a new user
	if err := handler.repo.Create(user); err != nil {
		handler.log.Error("Error creating user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
	}

	// Create JWT token with user ID, username, and role.
	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Username, user.Role, handler.config.JWTSecret)
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
	response := ToCreateUserResponse(user, accessToken, refreshToken)

	handler.log.Info("User created successfully", zap.String("userID", user.ID))
	return c.Status(fiber.StatusCreated).JSON(response)
}

// getEndpoint retrieves a user by their ID and returns their details.
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

	userResponse := ToResponseUser(user)

	handler.log.Info("User found:", zap.String("username", user.Username), zap.String("email", user.Email))

	return c.Status(fiber.StatusOK).JSON(userResponse)
}

// getAllEndpoint retrieves all users from the database.
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
		userResponses[i] = ToResponseUser(user)
	}

	handler.log.Info("All users fetched successfully")
	return c.Status(fiber.StatusOK).JSON(userResponses)
}

// updateEndpoint allows a user to update their own data if authorized.
func (handler *user) updateEndpoint(c *fiber.Ctx) error {
	// Receiving the user ID from the request parameters (URL).
	userID := c.Params("id")

	// Extracting the user ID from the JWT token (authenticated user).
	tokenUserID := c.Locals("user_id").(string)
	handler.log.Info("UserID param:", zap.String("userID", userID))

	// Checking if the user is authorized to update only their own data.
	if tokenUserID != userID {
		// If the user tries to update someone else's data, return an unauthorized response.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to update this user",
		})
	}

	// Parsing the update data from the request body.
	updateData := new(model.User)
	if err := c.BodyParser(updateData); err != nil {
		// If the request body is invalid, return a bad request response.
		handler.log.Error("Error parsing update data", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Attempting to update the user's data in the database.
	err := handler.repo.UpdateOneByID(userID, updateData)
	if err != nil {
		// If the update operation fails, return an internal server error response.
		handler.log.Error("Error updating user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
	}

	// Logging the success of the update operation.
	handler.log.Info("User updated successfully", zap.String("userID", userID))

	// Returning a success response after the update is completed.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user_id": userID,
	})
}

// deleteEndpoint allows a user to delete their own account if authorized.
func (handler *user) deleteEndpoint(c *fiber.Ctx) error {
	// Receiving the user ID from the request parameters (URL).
	userID := c.Params("id")

	// Extracting the user ID from the JWT token (authenticated user).
	tokenUserID := c.Locals("user_id").(string)

	// Checking if the user is authorized to delete only their own data.
	if tokenUserID != userID {
		// If the user tries to delete someone else's data, return an unauthorized response.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to delete this user",
		})
	}

	// Checking if the user exists in the database before attempting to delete.
	_, err := handler.repo.FindOneByID(userID)
	if err != nil {
		// If the user is not found, return a not found response.
		handler.log.Error("User not found", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Deleting the user from the database.
	if err := handler.repo.DeleteOneByID(userID); err != nil {
		// If the delete operation fails, return an internal server error response.
		handler.log.Error("Error deleting user", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Error deleting user")
	}

	// Logging the success of the delete operation.
	handler.log.Info("User deleted successfully", zap.String("userID", userID))

	// Returning a success response after the delete is completed.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
		"user_id": userID,
	})
}

// findByEmailEndpoint allows searching for a user by email.
func (handler *user) findByEmailEndpoint(c *fiber.Ctx) error {
	// Get query parameter for email
	email := c.Query("email")

	// Email parameter control
	if email == "" {
		handler.log.Error("Email query parameter is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email query parameter is required"})
	}

	// Validate email format using the validator
	if !handler.validate.ValidateEmailFormat(email) {
		handler.log.Error("Invalid email format", zap.String("email", email))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
	}

	// Verifying the email through logging
	handler.log.Info("Searching for user with email:", zap.String("email", email))

	// Search the user in the database using UserService
	user, err := handler.userService.FindByEmail(email)
	if err != nil {
		handler.log.Error("User not found by email", zap.String("email", email), zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// User found, create response
	userResponse := ToResponseUser(user)

	handler.log.Info("User found by email", zap.String("email", email))
	return c.Status(fiber.StatusOK).JSON(userResponse)
}
