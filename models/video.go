package models

import "time"

type Video struct {
	Id    uint
	Title string `gorm:"not null"`
	//视频文件存放路径
	Path string `gorm:"not null;unique"`
	//外键
	UserID    uint
	User      User
	CreatedAt time.Time
}
