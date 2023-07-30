/*
该文件下保存着基础接口的控制器
*/
package controllers

import (
	"app/database"
	"app/models"
	"app/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 注册
func Register(c *gin.Context) {
	u := new(models.User)
	c.Bind(u)
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
		c.Abort()
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
		c.Abort()
		return
	}
	//生成token
	t, _ := utils.GenerateToken(u.Username, u.Password, u.Id)
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "注册成功。",
		"user_id":     u.Id,
		"token":       t,
	})
}

// 登录
func Login(c *gin.Context) {
	u := new(models.User)
	c.Bind(u)
	//查询是否在数据库中
	u_in_database := new(models.User)
	database.Handler.Where("username = ?", u.Username).First(u_in_database)
	if u_in_database.Id == 0 {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  fmt.Sprintf("未找到名为%s的用户！", u.Username),
			"user_id":     0,
			"token":       "",
		})
		c.Abort()
		return
	}
	//检查密码是否一致
	if u_in_database.Password != u.Password {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "密码错误！请重试。",
			"user_id":     0,
			"token":       "",
		})
		c.Abort()
		return
	}
	//生成token
	t, _ := utils.GenerateToken(u.Username, u.Password, u.Id)
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "登录成功，欢迎您。",
		"user_id":     u.Id,
		"token":       t,
	})
}
