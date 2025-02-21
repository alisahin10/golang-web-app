package main

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/handlers"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/middleware"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/services"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/validator"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewDevelopment()
	errors := middleware.AppError{}
	// Initialize local repository
	localDbPath := os.Getenv("LOCAL_DB_PATH")
	localRepo, err := local.NewBuntRepository(localDbPath)
	if err != nil {
		logger.Fatal("Error creating brand-new bunt local repository", zap.String("local_db_path_env_variable", localDbPath), zap.Error(err))
	}
	defer localRepo.Close()

	// Initialize validator
	validate := validator.NewValidator() // No repository passed

	// Load JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	// Initialize AppConfig with the JWT secret
	config := &handlers.AppConfig{
		JWTSecret: []byte(jwtSecret),
	}

	// Initialize UserService
	userService := services.NewUserService(localRepo) // Create the UserService instance

	// Initialize fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			switch e := err.(type) {
			case *fiber.Error:
				return ctx.Status(e.Code).JSON(e)
			}
			return nil
		},
		AppName: "Golang Web Application",
	})

	// Initialize auth-handler and pass the config containing JWT secret
	authHandler := handlers.NewAuth(logger, localRepo, validate, config, errors)
	authHandler.AssignEndpoints("auth", app)

	// Initialize user-handler and pass the config containing JWT secret and userService
	userHandler := handlers.NewUser(logger, localRepo, validate, config, userService, errors)
	userHandler.AssignEndpoints("/user", app)

	// Start listening on port 8080
	if err = app.Listen(":8080"); err != nil {
		logger.Fatal("Application terminated with an error", zap.Error(err))
	}
}
