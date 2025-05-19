package models

import (
	"time"
	
	"github.com/google/uuid"
)

type Page struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key"`
	Title       string
	Content     string
	ParentID    *uuid.UUID `gorm:"type:uuid"`
	UserID      uuid.UUID  `gorm:"type:uuid"`
	IsDeleted   bool
	IsFavorite  bool
	Emoji       string
	Icon        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PageVersion struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	PageID    uuid.UUID `gorm:"type:uuid"`
	Content   string
	CreatedAt time.Time
}

type Discussion struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	PageID    uuid.UUID `gorm:"type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Content   string
	CreatedAt time.Time
}