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

	// Initialize database connection
	config.InitializeDB()

	// Create Gin router
	r := gin.Default()

	// Initialize handlers with database instance
	authHandler := handlers.NewAuthHandler(config.DB)

	// Public routes
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)
	r.POST("/reset-password", authHandler.RequestPasswordReset)

	// Protected admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminMiddleware(config.DB)) // Inject DB to admin middleware
	{
		admin.GET("/approve-user/:id", authHandler.ApproveUser)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}