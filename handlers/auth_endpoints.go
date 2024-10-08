package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/utils/jwt"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/validator"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AppConfig contains the JWT secret and other global configurations
type AppConfig struct {
	JWTSecret []byte // JWT secret used for signing tokens
}

type Auth struct {
	log      *zap.Logger        // Logger for logging events
	repo     local.Repository   // Repository interface for database operations
	validate validator.Validate // Validator for input validation
	config   *AppConfig         // Application configuration, including JWT secret
}

// NewAuth initializes a new Auth handler with its dependencies.
func NewAuth(log *zap.Logger, repo local.Repository, validate validator.Validate, config *AppConfig) Handler {
	return &Auth{
		log:      log,
		repo:     repo,
		validate: validate,
		config:   config,
	}
}

// AssignEndpoints sets up the routes for login, logout, and token refresh.
func (handler *Auth) AssignEndpoints(prefix string, router fiber.Router) {
	r := router.Group(prefix)

	// Route for user login
	r.Post("login", handler.loginEndpoint) // POST /auth/login: Authenticates user and returns tokens.

	// Route for user logout
	r.Post("logout", handler.logoutEndpoint) // POST /auth/logout: Logs the user out by invalidating their refresh token.

	// Route for refreshing tokens
	r.Post("refresh", handler.refreshTokenEndpoint) // POST /auth/refresh: Generates new access and refresh tokens.
}

// loginEndpoint handles the login process by validating the user, checking credentials, and generating tokens.
func (handler *Auth) loginEndpoint(ctx *fiber.Ctx) error {
	handler.log.Info("login endpoint called")

	// Struct for parsing login request body
	type loginRequest struct {
		Email    string `json:"email"`    // User's email address
		Password string `json:"password"` // User's password
	}
	var req loginRequest

	// Parse the request body into loginRequest struct
	if err := ctx.BodyParser(&req); err != nil {
		handler.log.Error("failed to parse body", zap.Error(err))
		return fiber.ErrBadRequest // Return 400 Bad Request if parsing fails
	}

	// Validate the parsed request (e.g., check if Email and Password are present)
	if err := handler.validate.Struct(req); err != nil {
		handler.log.Error("failed to validate body", zap.Error(err))
		return fiber.ErrBadRequest // Return 400 if validation fails
	}

	// Fetch the user from the repository using their email
	user, err := handler.repo.FindOneByEmail(req.Email)
	if err != nil {
		handler.log.Error("failed to find user by email", zap.Error(err))
		return &fiber.Error{Code: fiber.StatusUnauthorized, Message: "invalid username or password"} // This could be improved to return 404 if user not found
	}

	// Compare the provided password with the hashed password from the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		handler.log.Error("invalid password", zap.Error(err))
		return &fiber.Error{Code: fiber.StatusUnauthorized, Message: "invalid username or password"} // Return 401 Unauthorized if password is incorrect
	}

	// Delete any previous refresh token before issuing a new one
	if err := handler.repo.DeleteRefreshToken(user.ID); err != nil {
		if err.Error() == "not found" {
			handler.log.Info("No previous refresh token found", zap.String("userID", user.ID))
		} else {
			handler.log.Error("Failed to delete previous refresh token", zap.Error(err))
		}
	}

	// Generate access and refresh tokens using the JWT secret from the AppConfig
	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Username, user.Role, handler.config.JWTSecret)
	if err != nil {
		handler.log.Error("failed to generate tokens", zap.Error(err))
		return fiber.ErrInternalServerError // Return 500 if token generation fails
	}

	// Save the refresh token in the database
	err = handler.repo.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		handler.log.Error("failed to save refresh token", zap.Error(err))
		return fiber.ErrInternalServerError // Return 500 if refresh token save fails
	}

	// Respond with access token, refresh token, and user details
	return ctx.JSON(fiber.Map{
		"access_token":  accessToken,          // JWT access token for authentication
		"refresh_token": refreshToken,         // JWT refresh token for obtaining new access tokens
		"user":          ToResponseUser(user), // User details
	})
}

// logoutEndpoint handles user logout operations.
func (handler *Auth) logoutEndpoint(ctx *fiber.Ctx) error {
	// Parse request body to get the token
	type logoutRequest struct {
		Token string `json:"token"` // Refresh token to be invalidated
	}

	var req logoutRequest
	if err := ctx.BodyParser(&req); err != nil {
		handler.log.Error("Failed to parse body", zap.Error(err))
		return fiber.ErrBadRequest // Return 400 if body parsing fails
	}

	if req.Token == "" {
		return fiber.ErrBadRequest // Return 400 if no token is provided
	}

	handler.log.Info("Logout successful", zap.String("token", req.Token))

	// Find the user associated with the refresh token by its value
	userID, err := handler.repo.FindRefreshToken(req.Token)
	if err != nil {
		handler.log.Error("Invalid refresh token", zap.Error(err))
		return &fiber.Error{Code: fiber.StatusUnauthorized, Message: "You are not logged in"} // If token is invalid, respond with 401
	}

	// Delete the refresh token from the database by userID
	if err := handler.repo.DeleteRefreshToken(userID); err != nil {
		handler.log.Error("Failed to delete refresh token", zap.Error(err))
		return fiber.ErrInternalServerError // Return 500 if refresh token deletion fails
	}

	// Success response
	return ctx.JSON(fiber.Map{
		"message": "Logged out successfully", // Successful logout message
	})
}

// refreshTokenEndpoint creates new access and refresh tokens using the current refresh token.
func (handler *Auth) refreshTokenEndpoint(ctx *fiber.Ctx) error {
	type refreshRequest struct {
		Identifier   string `json:"identifier"`    // Can be username or email
		RefreshToken string `json:"refresh_token"` // The refresh token provided by the user
	}

	var req refreshRequest
	// Parse the request body to get the refresh token and identifier
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest // 400 - Bad request if the body is malformed
	}

	// Check if the identifier or token is empty
	if req.Identifier == "" || req.RefreshToken == "" {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "There is no refresh token\""} // 400 - Bad request if identifier or token is missing
	}

	// Verify if the refresh token has expired
	if jwt.IsExpired(req.RefreshToken, handler.config.JWTSecret) {
		return &fiber.Error{Code: fiber.StatusUnauthorized, Message: "Token has expired"} // 401 - Unauthorized if the token has expired
	}

	// Verify the refresh token by searching for it in the database
	userID, err := handler.repo.FindRefreshToken(req.RefreshToken)
	if err != nil {
		handler.log.Error("Invalid refresh token", zap.Error(err))
		return &fiber.Error{Code: fiber.StatusUnauthorized, Message: "Invalid refresh token"} // 401 - Unauthorized if token is not found or invalid
	}

	// Fetch the user using the identifier (could be email or username)
	user, err := handler.repo.FindOneByID(userID)
	if err != nil || (user.Email != req.Identifier && user.Username != req.Identifier) {
		handler.log.Error("User not found or identifier mismatch", zap.Error(err))
		return &fiber.Error{Code: fiber.StatusUnauthorized, Message: "There is no matching user"} // 401 - Unauthorized if the user is not found or identifier mismatch
	}

	// Delete any previous refresh token before issuing a new one
	if err := handler.repo.DeleteRefreshToken(user.ID); err != nil {
		handler.log.Error("Failed to delete previous refresh token", zap.Error(err))
	}

	// Generate new access and refresh tokens using the JWT secret from the AppConfig
	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Username, user.Role, handler.config.JWTSecret)
	if err != nil {
		handler.log.Error("Failed to generate tokens", zap.Error(err))
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Failed to generate tokens"}
		// 500 - Internal server error if token generation fails
	}

	// Save the new refresh token in the database
	if err := handler.repo.SaveRefreshToken(user.ID, refreshToken); err != nil {
		handler.log.Error("Failed to save new refresh token", zap.Error(err))
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Failed to save new refresh token"} // 500 - Internal server error if saving refresh token fails
	}

	// Use the ToCreateUserResponse function to generate the response
	response := ToCreateUserResponse(user, accessToken, refreshToken)

	// Respond with the structured response
	handler.log.Info("Successfully created new token", zap.String("user", req.Identifier))
	return ctx.Status(fiber.StatusOK).JSON(response) // Return successful response with new tokens
}
