package models

import (
	"gorm.io/gorm"
	"time"
)

type Page struct {
	gorm.Model
	Title      string `gorm:"not null"`
	Content    string `gorm:"type:text"`
	UserID     uint   `gorm:"not null"`
	IsDeleted  bool   `gorm:"default:false"`
	IsFavorite bool   `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}