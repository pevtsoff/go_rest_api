package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string
	Posts []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
	gorm.Model
	Title  string
	Body   string
	UserID uint
}
