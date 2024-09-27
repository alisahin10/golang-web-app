package main

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/handlers"
	"gitlab.com/rapsodoinc/tr/architecture/golang-web-app/repository/local"
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
	authHandler := handlers.NewAuth(log, localRepo)
	authHandler.AssignEndpoints("auth", app)

	// Start listening a port to be able to serve the http server
	if err = app.Listen(":8080"); err != nil {
		log.Fatal("Application terminated with an error", zap.Error(err))
	}
}
