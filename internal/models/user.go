//package models
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
/*
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
}*/
// File: internal/models/user.go (update User model)
//package models

/*import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusApproved UserStatus = "approved"
	UserStatusRejected UserStatus = "rejected"
)

type User struct {
	ID                uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Name              string      `gorm:"not null"`
	Email             string      `gorm:"unique;not null"`
	Password          string      `gorm:"not null"`
	Status            UserStatus  `gorm:"type:varchar(20);default:'pending'"`
	IsAdmin           bool        `gorm:"default:false"`
	CurrentWorkspaceID uuid.UUID  `gorm:"type:uuid"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}*/
//package models
/*
import (
	"time"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusApproved UserStatus = "approved"
	UserStatusRejected UserStatus = "rejected"
)

type User struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Name        string      `gorm:"not null"`
	Email       string      `gorm:"unique;not null"`
	Password    string      `gorm:"not null"`
	Status      UserStatus  `gorm:"type:varchar(20);default:'pending'"`
	IsAdmin     bool        `gorm:"default:false"`
	WorkspaceID uuid.UUID   `gorm:"type:uuid"` // Changed from CurrentWorkspaceID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	
	// Add relationship to workspaces
	Workspaces []*Workspace `gorm:"many2many:user_workspaces;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

// New model for workspace membership
type UserWorkspace struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	IsOwner     bool      `gorm:"default:false"`
	CreatedAt   time.Time
}

// Workspace model
type Workspace struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	
	// Add relationship to users
	Users []*User `gorm:"many2many:user_workspaces;"`
}

func (w *Workspace) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New()
	return
}*/
package models

import (
	"time"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusApproved UserStatus = "approved"
	UserStatusRejected UserStatus = "rejected"
)

type User struct {
	ID                uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Name              string      `gorm:"not null"`
	Email             string      `gorm:"unique;not null"`
	Password          string      `gorm:"not null"`
	Status            UserStatus  `gorm:"type:varchar(20);default:'pending'"`
	IsAdmin           bool        `gorm:"default:false"`
	CurrentWorkspaceID uuid.UUID  `gorm:"type:uuid"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}