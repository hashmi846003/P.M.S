package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Status    string `gorm:"default:'pending'"`
	Role      string `gorm:"default:'user'"`
	Pages     []Page `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}