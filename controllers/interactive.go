package controllers

import (
	"app/construct"
	"app/database"
	"app/models"
	"app/utils"
	"fmt"
	"strconv"
	"time"

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

// 评论操作
func Comment(c *gin.Context) {
	//获取必要参数
	tokenString := c.Query("token")
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	action_type, _ := strconv.Atoi(c.Query("action_type"))
	fmt.Println("测试")
	fmt.Println(video_id)
	fmt.Println(action_type)
	//先鉴权
	if !utils.CheckToken(tokenString) {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "token鉴定失败！",
			"comment":     nil,
		})
		c.Abort()
		return
	}
	//解析token
	claim, _ := utils.ParseToken(tokenString)
	if action_type == 1 {
		//获取评论内容
		comment_text := c.Query("comment_text")
		comment := models.Comment{
			Content:     comment_text,
			VideoID:     uint(video_id),
			UserID:      claim.Id,
			CreatedTime: time.Now().Unix(),
		}
		database.Handler.Create(&comment)
		c.JSON(200, gin.H{
			"status_code": 0,
			"status_msg":  "评论成功。",
			"comment":     comment,
		})
	} else {
		//获取需要删除的评论id
		comment_id, _ := strconv.Atoi(c.Query("comment_id"))
		//检查评论是否存在已经作者是否是本人
		comment := new(models.Comment)
		database.Handler.Where("id = ? AND user_id = ?", comment_id, claim.Id).First(comment)
		if comment.Id == 0 {
			c.JSON(200, gin.H{
				"status_code": 1,
				"status_msg":  "评论删除失败，评论不存在或者无权删除其他人的评论！",
				"comment":     nil,
			})
		} else {
			//删除操作
			database.Handler.Delete(comment)
			c.JSON(200, gin.H{
				"status_code": 0,
				"status_msg":  "评论删除成功。",
				"comment":     nil,
			})
		}
	}
}

// 获取评论列表
func GetCommentList(c *gin.Context) {
	//获取必要参数
	tokenString := c.Query("token")
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	//先鉴权
	if !utils.CheckToken(tokenString) {
		c.JSON(200, gin.H{
			"status_code":  1,
			"status_msg":   "token鉴定失败！",
			"comment_list": nil,
		})
		c.Abort()
		return
	}
	//解析token
	claim, _ := utils.ParseToken(tokenString)
	//查询所有评论
	var comments []models.Comment
	var comments_json []models.CommentJSON
	database.Handler.Where("video_id = ?", video_id).Find(&comments)
	for _, com := range comments {
		comment_json := construct.CommentJSON(com.Id, claim.Id)
		comments_json = append(comments_json, comment_json)
	}
	c.JSON(200, gin.H{
		"status_code":  0,
		"status_msg":   "查询评论列表成功。",
		"comment_list": comments_json,
	})
}
