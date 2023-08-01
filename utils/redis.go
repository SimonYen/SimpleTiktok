package utils

import (
	"app/database"
	"fmt"
)

// 获取视频点赞数量
func GetVideoLikeCount(video_id uint) int64 {
	key := fmt.Sprintf("like:v%d", video_id)
	res, err := database.RDB.SCard(database.CTX, key).Result()
	if err != nil {
		return 0
	}
	return res
}

// 视频赞操作
func LikeOrDisLikeVideo(video_id uint, user_id uint, action_type int) bool {
	key := fmt.Sprintf("like:v%d", video_id)
	member := fmt.Sprintf("u%d", user_id)
	if action_type == 1 {
		_, err := database.RDB.SAdd(database.CTX, key, member).Result()
		return err != nil
	} else {
		_, err := database.RDB.SRem(database.CTX, key, member).Result()
		return err != nil
	}
}

// 是否点赞
func VideoIsLiked(video_id, user_id uint) bool {
	key := fmt.Sprintf("like:v%d", video_id)
	member := fmt.Sprintf("u%d", user_id)
	res, _ := database.RDB.SIsMember(database.CTX, key, member).Result()
	return res
}
