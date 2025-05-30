package main

import (
	"log"
	"net/http"

	"github.com/Andrew-Ayman123/GoProject/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/files"

	_ "github.com/Andrew-Ayman123/GoProject/docs" // Import the generated docs

	"github.com/Andrew-Ayman123/GoProject/utils/env"
)

// @title Go Project API
// @version 1.0
// @description This is a sample server for a Go project.
// @termsOfService http://swagger.io/terms/

// @contact.name Andrew Ayman
// @contact.url https://github.com/Andrew-Ayman123
// @contact.email andrewayman9@gmail.com

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description API Key for authorization

func main() {
	env.Init() // Initialize environment variables

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	

	// Define API routes under /api/v1
	r.Route("/api/v1", func(v1 chi.Router) {

		// Swagger endpoint
		v1.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/api/v1/swagger/doc.json"),
		))

		// Group /user routes
		v1.Route("/user", func(user chi.Router) {
			user.Post("/login", handler.HandleUserLogIn)
		})
	})

	addr:= env.GetEnv("ADDR", ":8080")
	
	log.Default().Println("Starting server on localhost", addr)
	log.Default().Printf("Swagger UI available at: http://localhost%s/api/v1/swagger/index.html\n", addr)

	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
