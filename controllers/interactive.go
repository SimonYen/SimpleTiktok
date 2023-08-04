package controllers

import (
	"app/construct"
	"app/database"
	"app/models"
	"app/utils"
	"fmt"
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
	database.Handler.Where("user_id = ? AND id = ?", claim.Id, uint(video_id)).First(video)
	fmt.Println(video_id)
	if video.Id != 0 {
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

// 喜欢列表
func GetLikeList(c *gin.Context) {
	//获取必要参数
	tokenString := c.Query("token")
	user_id, _ := strconv.Atoi(c.Query("user_id"))
	//先鉴权
	if !utils.CheckToken(tokenString) {
		c.JSON(200, gin.H{
			"status_code": "1",
			"status_msg":  "token鉴定失败！",
			"video_list":  nil,
		})
		c.Abort()
		return
	}
	//解析token
	claim, _ := utils.ParseToken(tokenString)
	//获取用户喜欢的视频列表
	video_list := utils.GetUserLikedVideoList(uint(user_id))
	var videos []models.VideoJSON
	for _, v_id := range video_list {
		video_id, _ := strconv.Atoi(v_id)
		v := construct.VideoJSON(uint(video_id), claim.Id)
		videos = append(videos, v)
	}
	c.JSON(200, gin.H{
		"status_code": "0",
		"status_msg":  "获取喜欢列表成功。",
		"video_list":  videos,
	})
}
