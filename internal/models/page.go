/*package models

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	Title      string       `gorm:"size:255;index"`
	Content    string       `gorm:"type:text"`
	AuthorID   uint         `gorm:"index"`
	IsFavorite bool         `gorm:"default:false;index"`
	IsTrash    bool         `gorm:"default:false;index"`
	ParentID   *uint        `gorm:"index"`
	Children   []Page       `gorm:"foreignKey:ParentID"`
	Discussions []Discussion `gorm:"foreignKey:PageID"`
	Histories   []PageHistory `gorm:"foreignKey:PageID"`
}

type PageHistory struct {
	gorm.Model
	PageID    uint
	Content   string `gorm:"type:text"`
	Version   int
	UpdatedBy uint `gorm:"index"`
}*/
package models
import "gorm.io/gorm"
type Page struct {
	gorm.Model
	Title      string
	Content    string `gorm:"type:text"`
	AuthorID   uint
	IsFavorite bool
	IsTrash    bool
	ParentID   *uint
}