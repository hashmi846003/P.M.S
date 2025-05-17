package main

import (
	"os"
	"internal/config"
	"internal/handlers"
	"internal/middleware"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.InitializeDB()

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authHandler := handlers.NewAuthHandler(config.DB)
	pageHandler := handlers.NewPageHandler(config.DB)

	// Public routes
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)

	// Authenticated routes
	auth := r.Group("/")
	auth.Use(
		middleware.AuthMiddleware(),
		middleware.UserLoaderMiddleware(config.DB),
		middleware.ActiveUserMiddleware(),
	)
	{
		// Pages
		auth.GET("/pages", pageHandler.GetAllPages)
		auth.POST("/pages", pageHandler.CreatePage)
		auth.GET("/pages/:id", pageHandler.GetPage)
		auth.PUT("/pages/:id", pageHandler.UpdatePage)
		auth.DELETE("/pages/:id", pageHandler.DeletePage)
		auth.POST("/pages/:id/duplicate", pageHandler.DuplicatePage)
		auth.POST("/pages/:id/favorite", pageHandler.ToggleFavorite)
		auth.POST("/pages/:id/move-to-trash", pageHandler.MoveToTrash)
		auth.POST("/pages/:id/restore", pageHandler.RestorePage)

		// Discussions
		auth.POST("/pages/:id/discussions", pageHandler.AddDiscussion)

		// Search
		auth.GET("/search", pageHandler.SearchPages)
	}

	// Admin routes
	admin := r.Group("/admin")
	admin.Use(
		middleware.AdminMiddleware(config.DB),
		middleware.RateLimiter(50, time.Minute),
	)
	{
		admin.GET("/trash", pageHandler.GetTrash)
		admin.DELETE("/trash/:id", pageHandler.PermanentlyDeletePage)
	}

	r.Run(":" + os.Getenv("PORT"))
}