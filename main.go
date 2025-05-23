package main
/*
import (
	"log"
	"os"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/hashmi846003/P.M.S/internal/handlers"
	"github.com/hashmi846003/P.M.S/internal/middleware"
	"github.com/hashmi846003/P.M.S/internal/models"
	"github.com/hashmi846003/P.M.S/internal/repository"
	//"github.com/hashmi846003/P.M.S/internal"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file - make sure it exists in the project root")
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Run migrations
	if err := db.AutoMigrate(
		&models.User{},
		&models.Page{},
		&models.PageVersion{},
		&models.Discussion{},
	); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	pageRepo := repository.NewPageRepository(db)
	discussionRepo := repository.NewDiscussionRepository(db)

	// Initialize handlers and middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepo)
	pageHandler := handlers.NewPageHandler(pageRepo, discussionRepo)
	emojiHandler := handlers.NewEmojiHandler()

	// Create Gin router with middleware
	router := gin.Default()
	
	// Public routes
	public := router.Group("/auth")
	public.POST("/signup", authMiddleware.SignupHandler)
	public.POST("/login", authMiddleware.LoginHandler)

	// Authenticated routes
	auth := router.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// Page management
		auth.GET("/pages", pageHandler.ListPages)
		auth.POST("/pages", pageHandler.CreatePage)
		auth.GET("/pages/:id", pageHandler.GetPage)
		auth.PUT("/pages/:id", pageHandler.UpdatePage)
		auth.DELETE("/pages/:id", pageHandler.DeletePage)
		auth.POST("/pages/:id/restore", pageHandler.RestorePage)
		auth.POST("/pages/:id/favorite", pageHandler.ToggleFavorite)
		auth.POST("/pages/:id/duplicate", pageHandler.DuplicatePage)
		auth.GET("/pages/:id/versions", pageHandler.GetVersions)

		// Discussions
		auth.POST("/pages/:id/discussions", pageHandler.CreateDiscussion)
		auth.GET("/pages/:id/discussions", pageHandler.GetDiscussions)

		// Formatting operations
		auth.POST("/pages/:id/format", pageHandler.FormatContent)
		auth.POST("/pages/:id/align", pageHandler.AlignText)

		// Emoji operations
		auth.GET("/emojis", emojiHandler.ListEmojis)
		auth.GET("/emojis/categories", emojiHandler.GetCategories)
		auth.POST("/pages/:id/emoji", pageHandler.AddEmoji)

		// Sharing and utilities
		auth.POST("/pages/:id/share", pageHandler.GenerateShareLink)
		auth.GET("/trash", pageHandler.ListTrash)
		auth.POST("/pages/:id/move-to-trash", pageHandler.MoveToTrash)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Server running on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func initDatabase() (*gorm.DB, error) {
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
}*/
//package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/hashmi846003/P.M.S/internal/handlers"
	"github.com/hashmi846003/P.M.S/internal/middleware"
	"github.com/hashmi846003/P.M.S/internal/models"
	"github.com/hashmi846003/P.M.S/internal/repository"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file - make sure it exists in the project root")
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Run migrations with UUID extension
	if err := runMigrations(db); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	pageRepo := repository.NewPageRepository(db)
	discussionRepo := repository.NewDiscussionRepository(db)

	// Initialize handlers and middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepo)
	pageHandler := handlers.NewPageHandler(pageRepo, discussionRepo)
	emojiHandler := handlers.NewEmojiHandler()

	// Create Gin router with middleware
	router := gin.Default()
	
	// Public routes
	public := router.Group("/auth")
	public.POST("/signup", authMiddleware.SignupHandler)
	public.POST("/login", authMiddleware.LoginHandler)

	// Authenticated routes
	auth := router.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// Page management
		auth.GET("/pages", pageHandler.ListPages)
		auth.POST("/pages", pageHandler.CreatePage)
		auth.GET("/pages/:id", pageHandler.GetPage)
		auth.PUT("/pages/:id", pageHandler.UpdatePage)
		auth.DELETE("/pages/:id", pageHandler.DeletePage)
		auth.POST("/pages/:id/restore", pageHandler.RestorePage)
		auth.POST("/pages/:id/favorite", pageHandler.ToggleFavorite)
		auth.POST("/pages/:id/duplicate", pageHandler.DuplicatePage)
		auth.GET("/pages/:id/versions", pageHandler.GetVersions)

		// Discussions
		auth.POST("/pages/:id/discussions", pageHandler.CreateDiscussion)
		auth.GET("/pages/:id/discussions", pageHandler.GetDiscussions)

		// Formatting operations
		auth.POST("/pages/:id/format", pageHandler.FormatContent)
		auth.POST("/pages/:id/align", pageHandler.AlignText)

		// Emoji operations
		auth.GET("/emojis", emojiHandler.ListEmojis)
		auth.GET("/emojis/categories", emojiHandler.GetCategories)
		auth.POST("/pages/:id/emoji", pageHandler.AddEmoji)

		// Sharing and utilities
		auth.POST("/pages/:id/share", pageHandler.GenerateShareLink)
		auth.GET("/trash", pageHandler.ListTrash)
		auth.POST("/pages/:id/move-to-trash", pageHandler.MoveToTrash)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Server running on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func initDatabase() (*gorm.DB, error) {
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
}

func runMigrations(db *gorm.DB) error {
	// Create UUID extension first
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	// Run auto migrations
	return db.AutoMigrate(
		&models.User{},
		&models.Page{},
		&models.PageVersion{},
		&models.Discussion{},
		&models.ShareLink{},
	)
}