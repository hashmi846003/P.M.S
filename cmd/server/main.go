package main

import (
	"log"
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"your-app/internal/handlers"
	"your-app/internal/middleware"
	"your-app/internal/models"
	"your-app/internal/repository"
	"your-app/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(&models.User{}, &models.Page{}); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	pageRepo := repository.NewPageRepository(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo)
	pageHandler := handlers.NewPageHandler(pageRepo)

	// Create Gin router
	router := gin.Default()

	// Public routes
	router.POST("/signup", authHandler.Signup)
	router.POST("/login", authHandler.Login)

	// Authenticated routes
	authGroup := router.Group("/")
	authGroup.Use(middleware.JWTAuthMiddleware())
	{
		authGroup.GET("/pages", pageHandler.GetAllPages)
		authGroup.POST("/pages", pageHandler.CreatePage)
		authGroup.GET("/pages/:id", pageHandler.GetPageByID)
		authGroup.PUT("/pages/:id", pageHandler.UpdatePage)
		authGroup.DELETE("/pages/:id", pageHandler.DeletePage)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}