package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Andrew-Ayman123/GoProject/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/Andrew-Ayman123/GoProject/docs" // Import the generated docs package
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
	r := gin.Default()
	
	// API v1 group
	v1 := r.Group("/api/v1")

	userAuth := v1.Group("/user")
	{
		// Add your login route here
		userAuth.POST("/login", ginHandlerWrapper(handler.HandleUserLogIn))
		
	}
	
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Starting server on :8080")
	fmt.Println("Swagger UI available at: http://localhost:8080/api/v1/swagger/index.html")
	
	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

// Wrapper function to convert http.HandlerFunc to gin.HandlerFunc
func ginHandlerWrapper(h http.HandlerFunc) gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(h))
}