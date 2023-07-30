/*
该文件下保存着基础接口的控制器
*/
package controllers

import (
	"app/database"
	"app/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 注册
func Register(c *gin.Context) {
	u := new(models.User)
	//获取Query参数
	err := c.Bind(u)
	if err != nil {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "解析query参数失败！",
			"user_id":     0,
			"token":       "",
		})
		return
	}
	//检查是否为空
	if u.Username == "" || u.Password == "" {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "用户名或者密码不能为空！",
			"user_id":     0,
			"token":       "",
		})
		return
	}
	//查询是否有同名的用户
	u_in_database := new(models.User)
	database.Handler.Where("username = ?", u.Username).First(u_in_database)
	if u_in_database.Id != 0 {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  fmt.Sprintf("用户名%s已被注册！", u.Username),
			"user_id":     0,
			"token":       "",
		})
		return
	}
	//写入数据库
	res := database.Handler.Create(u)
	if res.Error != nil {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "数据库插入失败！",
			"user_id":     0,
			"token":       "",
		})
		return
	}
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "注册成功。",
		"user_id":     u.Id,
		"token":       "string",
	})
}
