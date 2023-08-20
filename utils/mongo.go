package utils

import (
	"app/database"
	"app/models"
	"context"
	"time"
)

// 将需要发送的消息存入数据库
func SaveMessage(from_user_id, to_user_id int, content string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//构造结构体
	msg := models.Message{
		FromUserID: int64(from_user_id),
		ToUserID:   int64(to_user_id),
		Content:    content,
		CreateTime: time.Now().Unix(),
	}
	collection := database.MON.Database("tiktok").Collection("messages")
	_, err := collection.InsertOne(ctx, msg)
	return err == nil
}
