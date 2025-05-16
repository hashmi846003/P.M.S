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
	Children   []Page      `gorm:"foreignkey:ParentID"`
	Discussions []Discussion
}

type Discussion struct {
	gorm.Model
	PageID  uint
	UserID  uint
	Content string `gorm:"type:text"`
}

type PageHistory struct {
	gorm.Model
	PageID    uint
	Content   string `gorm:"type:text"`
	Version   int
	UpdatedBy uint
}