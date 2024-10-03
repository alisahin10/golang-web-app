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

// NewAuth is the constructor for the auth handler. It initializes dependencies (logger, repo, validator)
func NewAuth(log *zap.Logger, repo local.Repository, validate validator.Validate) Handler {
	return &auth{
		log:      log,
		repo:     repo,
		validate: validate,
	}
}

// AssignEndpoints assigns routes for login and logout to the Fiber router.
func (handler *auth) AssignEndpoints(prefix string, router fiber.Router) {
	r := router.Group(prefix)

	r.Post("login", handler.loginEndpoint)
	r.Post("logout", handler.logoutEndpoint)
	r.Post("refresh", handler.refreshTokenEndpoint)

}

// loginEndpoint handles the login process by validating the user, checking credentials, and generating tokens.
func (handler *auth) loginEndpoint(ctx *fiber.Ctx) error {
	handler.log.Info("login endpoint called")

	// Struct for parsing login request body
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
		return fiber.ErrInternalServerError // This could be improved to return 404 if user not found
	}

	// Compare the provided password with the hashed password from the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		handler.log.Error("invalid password", zap.Error(err))
		return fiber.ErrUnauthorized // Change to Unauthorized (401) instead of Internal Server Error (500)
	}

	// Delete any previous refresh token before issuing a new one
	if err := handler.repo.DeleteRefreshToken(user.ID); err != nil {
		handler.log.Error("Failed to delete previous refresh token", zap.Error(err))
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.Username, user.Email)
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
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          utils.ToResponseUser(user),
	})

}

func (handler *auth) logoutEndpoint(ctx *fiber.Ctx) error {
	// Parse request body to get the token
	type logoutRequest struct {
		Token string `json:"token"`
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
		return fiber.ErrUnauthorized // If token is invalid, respond with 401
	}

	// Delete the refresh token from the database by userID
	if err := handler.repo.DeleteRefreshToken(userID); err != nil {
		handler.log.Error("Failed to delete refresh token", zap.Error(err))
		return fiber.ErrInternalServerError // Return 500 if refresh token deletion fails
	}

	// Success response
	return ctx.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (handler *auth) refreshTokenEndpoint(ctx *fiber.Ctx) error {
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
		return fiber.ErrBadRequest // 400 - Bad request if identifier or token is missing
	}

	// Verify if the refresh token has expired
	if utils.IsExpired(req.RefreshToken) {
		return fiber.ErrUnauthorized // 401 - Unauthorized if the token has expired
	}

	// Verify the refresh token by searching for it in the database
	userID, err := handler.repo.FindRefreshToken(req.RefreshToken)
	if err != nil {
		handler.log.Error("Invalid refresh token", zap.Error(err))
		return fiber.ErrUnauthorized // 401 - Unauthorized if token is not found or invalid
	}

	// Fetch the user using the identifier (could be email or username)
	user, err := handler.repo.FindOneByID(userID)
	if err != nil || (user.Email != req.Identifier && user.Username != req.Identifier) {
		handler.log.Error("User not found or identifier mismatch", zap.Error(err))
		return fiber.ErrUnauthorized // 401 - Unauthorized if the user is not found or identifier mismatch
	}

	// Generate new access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.Username, user.Email)
	if err != nil {
		handler.log.Error("Failed to generate tokens", zap.Error(err))
		return fiber.ErrInternalServerError // 500 - Internal server error if token generation fails
	}

	// Optionally, save the new refresh token in the database
	if err := handler.repo.SaveRefreshToken(user.ID, refreshToken); err != nil {
		handler.log.Error("Failed to save new refresh token", zap.Error(err))
		return fiber.ErrInternalServerError // 500 - Internal server error if saving refresh token fails
	}

	// Use the ToCreateUserResponse function to generate the response
	response := utils.ToCreateUserResponse(user, accessToken, refreshToken)

	// Respond with the structured response
	return ctx.Status(fiber.StatusOK).JSON(response)
}
