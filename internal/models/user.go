package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email         string `gorm:"unique;not null"`
	Password      string `gorm:"not null"`
	AdminApproval bool   `gorm:"default:false"`
	IsAdmin       bool   `gorm:"default:false"`
}