package models

type Comment struct {
	Id      uint
	Content string
	//外键
	VideoID uint
	Video   Video
	//再附加一个作者id，方便些
	UserID      uint
	CreatedTime int64 `gorm:"not null"`
}
