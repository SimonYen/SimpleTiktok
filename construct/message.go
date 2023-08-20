package construct

import (
	"app/database"
	"app/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// 直接生成结构体数组
func Messages(from_user_id, to_user_id int) []models.Message {
	collection := database.MON.Database("tiktok").Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cond := bson.M{
		"to_user_id":   to_user_id,
		"from_user_id": from_user_id,
	}
	cursor, err := collection.Find(ctx, cond)
	if err != nil {
		return nil
	}
	var ret []models.Message
	cursor.All(ctx, &ret)
	return ret
}
