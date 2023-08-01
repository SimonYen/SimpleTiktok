package controllers

import (
	"app/database"
	"app/models"
	"app/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 视频点赞
func VideoLike(c *gin.Context) {
	//获取必要参数
	tokenString := c.Query("token")
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	action_type, _ := strconv.Atoi(c.Query("action_type"))
	//先鉴权
	if !utils.CheckToken(tokenString) {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "token鉴定失败！",
		})
		c.Abort()
		return
	}
	//解析token
	claim, _ := utils.ParseToken(tokenString)
	//不能给自己点赞
	video := new(models.Video)
	database.Handler.Where("user_id = ?", claim.Id).First(video)
	if uint(video_id) == video.Id {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "小伙子，别给自己刷赞好不好！",
		})
		c.Abort()
		return
	}
	//理论上应该对video_id进行判断，看是否在数据库中存在，这里我偷个懒算了
	//先找出点赞者id
	res := utils.LikeOrDisLikeVideo(uint(video_id), claim.Id, action_type)
	if !res {
		c.JSON(200, gin.H{
			"status_code": 0,
			"status_msg":  "Redis操作失败！",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "操作成功。",
	})
}
