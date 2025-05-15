package main

import (
	"os"
	"internal/config"
	"internal/handlers"
	"internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize database
	config.InitializeDB()

	// Create Gin router
	r := gin.Default()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(config.DB)

	// Routes
	api := r.Group("/api/v1")
	{
		api.POST("/signup", authHandler.Signup)
		api.POST("/login", authHandler.Login)
		api.POST("/reset-password", authHandler.RequestPasswordReset)
		
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			admin.GET("/approve-user/:id", authHandler.ApproveUser)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}