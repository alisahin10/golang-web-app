package main

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/handlers"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/validator"
	"go.uber.org/zap"
	"os"
)

func main() {
	// Initialize logger
	log, _ := zap.NewDevelopment()

	// Initialize local repository
	localDbPath := os.Getenv("LOCAL_DB_PATH")
	localRepo, err := local.NewBuntRepository(localDbPath)
	if err != nil {
		log.Fatal("Error creating brand-new bunt local repository", zap.String("local_db_path_env_variable", localDbPath), zap.Error(err))
	}

	defer func() {
		if buntRepo, ok := localRepo.(*local.BuntImpl); ok {
			err := buntRepo.DB.Close()
			if err != nil {
				return
			} // Close the DB
		}
	}()

	// Initialize validator
	validate := validator.NewValidator(localRepo)

	// Initialize fiber app
	app := fiber.New(fiber.Config{
		AppName: "Golang Web Application",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Type your error handler here
			// But make your decisions wisely about this should be located right in the main method (!)
			return nil
		},
	})

	// Initialize auth-handler
	authHandler := handlers.NewAuth(log, localRepo, validate)
	authHandler.AssignEndpoints("auth", app)

	// Initialize user-handler
	userHandler := handlers.NewUser(log, localRepo, validate)
	userHandler.AssignEndpoints("/user", app)

	// Start listening a port to be able to serve the http server
	if err = app.Listen(":8080"); err != nil {
		log.Fatal("Application terminated with an error", zap.Error(err))
	}

}
