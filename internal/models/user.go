/*package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex;not null;size:255"`
	Password      string `gorm:"not null;size:255"`
	AdminApproval bool   `gorm:"default:false;index"`
	IsAdmin       bool   `gorm:"default:false;index"`
	Pages         []Page `gorm:"foreignKey:AuthorID"`
}*/
/*// user.go
package models

type User struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex;size:255"`
	Password      string `gorm:"size:255"`
	AdminApproval bool
	IsAdmin       bool
}

// page.go
package models

type Page struct {
	gorm.Model
	Title      string
	Content    string `gorm:"type:text"`
	AuthorID   uint
	IsFavorite bool
	IsTrash    bool
	ParentID   *uint
}*/
package models
import "gorm.io/gorm"
type User struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex;size:255"`
	Password      string `gorm:"size:255"`
	AdminApproval bool
	IsAdmin       bool
}