package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/model"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// login handles user login operations.
func (auth *Auth) login(ctx *fiber.Ctx) error {
	auth.log.Info("Login endpoint called")

	var req model.LoginRequest
	// Parse the request body
	if err := ctx.BodyParser(&req); err != nil {
		auth.log.Error("Error parsing request body", zap.Error(err))
		return fiber.ErrBadRequest // Return BadRequest on body parsing error
	}

	// Validation check
	if err := auth.validate.Struct(&req); err != nil {
		// Log validation errors
		for _, err := range err.(validator.ValidationErrors) {
			auth.log.Error("Validation failed",
				zap.String("field", err.StructNamespace()),
				zap.String("tag", err.Tag()),
				zap.String("value", err.Param()),
			)
		}
		return fiber.ErrBadRequest // Return BadRequest on validation error
	}

	// Find user by email
	user, err := auth.repo.FindOneByEmail(req.Email)
	if err != nil {
		auth.log.Error("Find user by email failed", zap.Error(err), zap.String("email", req.Email))
		return fiber.ErrUnauthorized // Return Unauthorized if user not found
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		auth.log.Error("Incorrect password", zap.String("email", req.Email))
		return fiber.ErrUnauthorized // Return Unauthorized if password is incorrect
	}

	// Generate tokens with JWT secret
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Username, user.Role, auth.config.JWTSecret)
	if err != nil {
		auth.log.Error("Generate access token failed", zap.Error(err))
		return fiber.ErrInternalServerError // Return InternalServerError on token generation failure
	}

	// Save refresh token
	if err := auth.repo.SaveRefreshToken(user.ID, refreshToken); err != nil {
		auth.log.Error("Save refresh token failed", zap.Error(err))
		return fiber.ErrInternalServerError // Return InternalServerError on refresh token save failure
	}

	// Return successful login response
	return ctx.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          utils.ToResponseUser(user), // Response with user information
	})
}

// logout handles user logout operations.
func (auth *Auth) logout(ctx *fiber.Ctx) error {
	auth.log.Info("Logout endpoint called")

	var req model.LogoutRequest
	// Parse the request body
	if err := ctx.BodyParser(&req); err != nil {
		auth.log.Error("Error parsing request body", zap.Error(err))
		return fiber.ErrBadRequest // Return BadRequest on body parsing error
	}

	// Token check
	if req.Token == "" {
		return ctx.JSON(fiber.Map{
			"status":  "logout unsuccessful", // Logout operation unsuccessful
			"message": "Token is required",   // Required token is missing
		})
	}

	// Find refresh token
	userID, err := auth.repo.FindRefreshToken(req.Token)
	if err != nil {
		auth.log.Error("Find refresh token failed", zap.Error(err))
		return fiber.ErrUnauthorized // Return Unauthorized if token not found
	}

	// Delete refresh token
	if err := auth.repo.DeleteRefreshToken(userID); err != nil {
		auth.log.Error("Delete refresh token failed", zap.Error(err))
		return fiber.ErrInternalServerError // Return InternalServerError on refresh token deletion failure
	}

	// Return successful response
	return ctx.JSON(fiber.Map{
		"status":  "logged out",              // Logout operation successful
		"message": "Logged out successfully", // Successful logout message
	})
}

// refreshToken creates new access and refresh tokens using the current refresh token.
func (auth *Auth) refreshToken(ctx *fiber.Ctx) error {
	auth.log.Info("Refresh token endpoint called")
	var req model.RefreshRequest
	// Parse the request body
	if err := ctx.BodyParser(&req); err != nil {
		auth.log.Error("Validation failed")
		return fiber.ErrBadRequest // Return BadRequest on body parsing error
	}
	// Return BadRequest if Identifier or RefreshToken is missing
	if req.Identifier == "" || req.RefreshToken == "" {
		return fiber.ErrBadRequest
	}
	// Check if the refresh token has expired using JWT secret
	if utils.IsExpired(req.RefreshToken, auth.config.JWTSecret) {
		return fiber.ErrBadRequest // Return BadRequest if expired
	}
	// Find user ID using refresh token
	userID, err := auth.repo.FindRefreshToken(req.RefreshToken)
	if err != nil {
		auth.log.Error("Find refresh token failed", zap.Error(err))
		return fiber.ErrInternalServerError // Return InternalServerError on refresh token retrieval failure
	}
	// Find user
	user, err := auth.repo.FindOneByEmail(userID)
	if err != nil || (user.Email != req.Identifier && user.Username != req.Identifier) {
		auth.log.Error("Find user by email failed", zap.Error(err))
		return fiber.ErrUnauthorized // Return Unauthorized if user not found or validation fails
	}
	// Delete old refresh token
	if err := auth.repo.DeleteRefreshToken(user.ID); err != nil {
		auth.log.Error("Delete refresh token failed", zap.Error(err))
	}
	// Generate new tokens with JWT secret
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Username, user.Role, auth.config.JWTSecret)
	if err != nil {
		auth.log.Error("Generate access token failed", zap.Error(err))
		return fiber.ErrInternalServerError // Return InternalServerError on token generation failure
	}
	// Save new refresh token
	if err := auth.repo.SaveRefreshToken(user.ID, refreshToken); err != nil {
		auth.log.Error("Save refresh token failed", zap.Error(err))
		return fiber.ErrInternalServerError // Return InternalServerError on new refresh token save failure
	}
	// Create successful response
	response := utils.ToCreateUserResponse(user, accessToken, refreshToken)
	auth.log.Info("Successfully created new token", zap.String("user", req.Identifier))
	return ctx.Status(fiber.StatusOK).JSON(response) // Return successful response with new tokens
}
