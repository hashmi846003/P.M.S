/*
package models

import (
	"time"
	"github.com/google/uuid"
)

type Permission string

const (
	PermissionViewer Permission = "viewer"
	PermissionEditor Permission = "editor"
)

type ShareLink struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	PageID     uuid.UUID  `gorm:"type:uuid;not null"`
	CreatedBy  uuid.UUID  `gorm:"type:uuid;not null"`
	Token      string     `gorm:"unique;not null"`
	Permission Permission `gorm:"type:varchar(10);default:'viewer'"`
	ExpiresAt  time.Time
	CreatedAt  time.Time
}*/
package models

import (
	"time"
	
	"github.com/google/uuid"
)

type Permission string

const (
	PermissionViewer Permission = "viewer"
	PermissionEditor Permission = "editor"
)

type ShareLink struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	PageID     uuid.UUID  `gorm:"type:uuid;not null"`
	CreatedBy  uuid.UUID  `gorm:"type:uuid;not null"`
	Token      string     `gorm:"unique;not null"`
	Permission Permission `gorm:"type:varchar(10);default:'viewer'"`
	ExpiresAt  time.Time
	CreatedAt  time.Time
}