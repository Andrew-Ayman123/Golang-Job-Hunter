package main

import (
	"log"

	"github.com/Andrew-Ayman123/Job-Hunter/internal/app"
	// "github.com/Andrew-Ayman123/Job-Hunter/utils/env"
)

// @title Go Project API
// @version 1.0
// @description This is a sample server for a Go project with JWT authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name Andrew Ayman
// @contact.url https://github.com/Andrew-Ayman123
// @contact.email andrewayman9@gmail.com

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using Bearer scheme. Example: "Bearer {token}"

func main() {
	// Initialize environment variables
	// env.Init()

	// Create and start the application
	application, err := app.NewApp()
	if err != nil {
		log.Fatal("Failed to create application: ", err)
	}

	if err := application.Run(); err != nil {
		log.Fatal("Failed to run application: ", err)
	}
}
