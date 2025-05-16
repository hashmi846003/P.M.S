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
	godotenv.Load()
	config.InitializeDB()

	r := gin.Default()

	authHandler := handlers.NewAuthHandler(config.DB)
	pageHandler := handlers.NewPageHandler(config.DB)

	// Public routes
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)

	// Authenticated routes
	auth := r.Group("/").Use(middleware.AuthMiddleware())
	{
		auth.GET("/pages", pageHandler.GetAllPages)
		auth.POST("/pages", pageHandler.CreatePage)
		auth.GET("/pages/:id", pageHandler.GetPage)
		auth.PUT("/pages/:id", pageHandler.UpdatePage)
		auth.DELETE("/pages/:id", pageHandler.DeletePage)
		auth.POST("/pages/:id/duplicate", pageHandler.DuplicatePage)
		auth.POST("/pages/:id/favorite", pageHandler.ToggleFavorite)
		auth.POST("/pages/:id/move-to-trash", pageHandler.MoveToTrash)
		auth.POST("/pages/:id/restore", pageHandler.RestorePage)
		auth.POST("/pages/:id/discussions", pageHandler.AddDiscussion)
		auth.GET("/search", pageHandler.SearchPages)
	}

	// Admin routes
	admin := auth.Group("/admin").Use(middleware.AdminMiddleware(config.DB))
	{
		admin.GET("/trash", pageHandler.GetTrash)
		admin.DELETE("/trash/:id", pageHandler.PermanentlyDeletePage)
	}

	r.Run(":" + os.Getenv("PORT"))
}