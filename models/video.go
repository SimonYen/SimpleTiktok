package models

type Video struct {
	Id    uint
	Title string `gorm:"not null"`
	//视频文件存放路径
	Path string `gorm:"not null;unique"`
	//外键
	UserID      uint
	User        User
	CreatedTime int64 `gorm:"not null"` //这个傻逼gorm框架，时间这个点坑了老子好久，我直接用时间戳了，操！
}

// 便于返回JSON数据的Video模型
type VideoJSON struct {
	Author        UserJSON `json:"author"`         // 视频作者信息
	CommentCount  int64    `json:"comment_count"`  // 视频的评论总数
	CoverURL      string   `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64    `json:"favorite_count"` // 视频的点赞总数
	ID            uint     `json:"id"`             // 视频唯一标识
	IsFavorite    bool     `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string   `json:"play_url"`       // 视频播放地址
	Title         string   `json:"title"`          // 视频标题
}
