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

type CommentJSON struct {
	Content    string   `json:"content"`     // 评论内容
	CreateDate string   `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64    `json:"id"`          // 评论id
	User       UserJSON `json:"user"`        // 评论用户信息
}
