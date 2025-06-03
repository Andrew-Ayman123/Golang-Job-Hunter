package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Andrew-Ayman123/Job-Hunter/internal/handlers"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/repository"
	routes "github.com/Andrew-Ayman123/Job-Hunter/internal/router"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/services"
	"github.com/Andrew-Ayman123/Job-Hunter/utils/env"
)

type App struct {
	server      *Server
	db          *sql.DB
	userRepo    repository.UserRepository
	jwtService  *services.JWTService
	userHandler *handlers.UserHandler
}

func NewApp() (*App, error) {
	// Initialize database
	db, err := newDatabaseConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize server
	addr := env.GetEnv("ADDR", ":8080")

	app := App{}

	// Initialize services
	app.jwtService = services.NewJWTService()

	// Initialize repositories
	app.userRepo = repository.NewUserRepository(db)

	// Initialize handlers
	app.userHandler = handlers.NewUserHandler(app.userRepo, app.jwtService)

	app.server = NewServer(addr, routes.SetupRoutes(app.userHandler, app.jwtService))

	return &app, nil
}

func (a *App) Run() error {
	log.Printf("Starting server on %s", a.server.Addr())
	log.Printf("Swagger UI available at: http://localhost%s/api/v1/swagger/index.html", a.server.Addr())

	return a.server.Start()
}

func (a *App) Shutdown() error {
	if err := a.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	if err := a.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
