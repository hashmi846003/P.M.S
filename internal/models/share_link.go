package models

import (
	"time"
	"github.com/google/uuid"
)

type ShareLink struct {
	Token     string    `gorm:"primaryKey;size:255"`
	PageID    uuid.UUID `gorm:"type:uuid"`
	ExpiresAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
}