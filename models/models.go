package models


import "gorm.io/gorm"

type User struct{
	gorm.Model
	Name string
}

type Post struct{
	gorm.Model
	Title string
	Body string
}