package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"your-app-name/internal/models"
)

var DB *gorm.DB

func InitializeDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	// Auto migrate models
	DB.AutoMigrate(&models.User{})
}