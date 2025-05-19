package models
/*
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
}*/
//package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string    `gorm:"unique"`
	Password  string
	Name      string
	CreatedAt time.Time
}

// Add this hook if using GORM
func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return nil
}