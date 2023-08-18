package controllers

import (
	"app/construct"
	"app/models"
	"app/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 关注操作
func Follow(c *gin.Context) {
	//获取必要参数
	tokenString := c.Query("token")
	//被关注的人的id
	to_user_id, _ := strconv.Atoi(c.Query("to_user_id"))
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
	claim, _ := utils.ParseToken(tokenString)
	res := utils.FollowOrUnfollowUser(claim.Id, uint(to_user_id), action_type)
	if res {
		c.JSON(200, gin.H{
			"status_code": 0,
			"status_msg":  "操作成功。",
		})
	} else {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "操作失败！",
		})
	}
}

// 关注列表
func GetFollowList(c *gin.Context) {
	//获取必要参数
	tokenString := c.Query("token")
	user_id, _ := strconv.Atoi(c.Query("user_id"))
	//先鉴权
	if !utils.CheckToken(tokenString) {
		c.JSON(200, gin.H{
			"status_code": "1",
			"status_msg":  "token鉴定失败！",
			"user_list":   nil,
		})
		c.Abort()
		return
	}
	claim, _ := utils.ParseToken(tokenString)
	//获取关注的用户名id
	res := utils.GetFollowList(uint(user_id))
	user_json_list := make([]models.UserJSON, 0, 30)
	for _, u_id_str := range res {
		u_id, _ := strconv.Atoi(u_id_str)
		tmp := construct.UserJSON(uint(u_id), claim.Id)
		user_json_list = append(user_json_list, tmp)
	}
	c.JSON(200, gin.H{
		"status_code": "0",
		"status_msg":  "获取关注列表成功。",
		"user_list":   user_json_list,
	})
}

// 粉丝列表
func GetFansList(c *gin.Context) {
	//获取必要参数
	tokenString := c.Query("token")
	user_id, _ := strconv.Atoi(c.Query("user_id"))
	//先鉴权
	if !utils.CheckToken(tokenString) {
		c.JSON(200, gin.H{
			"status_code": "1",
			"status_msg":  "token鉴定失败！",
			"user_list":   nil,
		})
		c.Abort()
		return
	}
	claim, _ := utils.ParseToken(tokenString)
	//获取粉丝的用户名id
	res := utils.GetFollowedList(uint(user_id))
	user_json_list := make([]models.UserJSON, 0, 30)
	for _, u_id_str := range res {
		u_id, _ := strconv.Atoi(u_id_str)
		tmp := construct.UserJSON(uint(u_id), claim.Id)
		user_json_list = append(user_json_list, tmp)
	}
	c.JSON(200, gin.H{
		"status_code": "0",
		"status_msg":  "获取粉丝列表成功。",
		"user_list":   user_json_list,
	})
}

// 好友列表
func GetFridendList(c *gin.Context) {
	//获取必要参数
	tokenString := c.Query("token")
	user_id, _ := strconv.Atoi(c.Query("user_id"))
	//先鉴权
	if !utils.CheckToken(tokenString) {
		c.JSON(200, gin.H{
			"status_code": "1",
			"status_msg":  "token鉴定失败！",
			"user_list":   nil,
		})
		c.Abort()
		return
	}
	claim, _ := utils.ParseToken(tokenString)
	//获取好友的用户名id
	res := utils.GetFriendList(uint(user_id))
	user_json_list := make([]models.UserJSON, 0, 30)
	for _, u_id_str := range res {
		u_id, _ := strconv.Atoi(u_id_str)
		tmp := construct.UserJSON(uint(u_id), claim.Id)
		user_json_list = append(user_json_list, tmp)
	}
	c.JSON(200, gin.H{
		"status_code": "0",
		"status_msg":  "获取好友列表成功。",
		"user_list":   user_json_list,
	})
}
