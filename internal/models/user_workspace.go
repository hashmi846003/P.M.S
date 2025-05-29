package models

import (
	"time"
	
	"github.com/google/uuid"
)

type UserWorkspace struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	IsOwner     bool      `gorm:"default:false"`
	CreatedAt   time.Time
}