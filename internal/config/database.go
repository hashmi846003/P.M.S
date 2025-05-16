package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"internal/models"
)

var DB *gorm.DB

func InitializeDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	DB.AutoMigrate(
		&models.User{},
		&models.Page{},
		&models.Discussion{},
		&models.PageHistory{},
	)
}