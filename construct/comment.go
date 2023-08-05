package construct

import (
	"app/database"
	"app/models"
	"time"
)

func CommentJSON(comment_id, operator_id uint) models.CommentJSON {
	//找出发布该评论的作者
	c := new(models.Comment)
	database.Handler.Where("id = ?", comment_id).First(c)
	//转换时间
	t := time.Unix(c.CreatedTime, 0).Format("01-02")
	return models.CommentJSON{
		ID:         int64(c.Id),
		User:       UserJSON(c.UserID, operator_id),
		Content:    c.Content,
		CreateDate: t,
	}
}
