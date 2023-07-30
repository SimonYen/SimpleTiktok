package models

type User struct {
	Id       uint
	Username string `form:"username" gorm:"unique;not null"`
	Password string `form:"password" gorm:"not null"`
}
