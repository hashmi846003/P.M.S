package config

import (
	"os"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"internal/models"
)

var DB *gorm.DB

func InitializeDB() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "test.db" // default for development
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Auto migrate models
	DB.AutoMigrate(
		&models.User{},
		&models.Page{},
		&models.Discussion{},
		&models.PageHistory{},
	)
}