package config

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitializeDB sets up the database connection based on environment configuration
func InitializeDB() {
	var dialector gorm.Dialector
	var err error

	dbType := os.Getenv("DB_TYPE")

	switch dbType {
	case "mysql":
		dialector = mysql.Open(getMySQLDSN())
	case "postgres":
		dialector = postgres.Open(getPostgresDSN())
	default: // Default to SQLite for local development
		dialector = sqlite.Open(getSQLitePath())
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	configureConnectionPool()
	autoMigrateModels()
}

func getSQLitePath() string {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "file:collab-doc.db?cache=shared&_fk=1" // Default SQLite with foreign key support
	}
	return dbPath
}

func getMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func getPostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
}

func configureConnectionPool() {
	sqlDB, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get database instance: %v", err))
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func autoMigrateModels() {
	err := DB.AutoMigrate(
		// List all your models here
		&User{},
		&Page{},
		&Discussion{},
		&PageHistory{},
		&Favorite{},
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to auto-migrate models: %v", err))
	}
}