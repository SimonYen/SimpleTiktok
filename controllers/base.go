/*
该文件下保存着基础接口的控制器
*/
package controllers

import (
	"app/config"
	"app/database"
	"app/models"
	"app/utils"
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	//生成随机渐变背景图
	utils.GenerateBackground(u.Id)
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
	t, _ := utils.GenerateToken(u.Username, u.Password, u_in_database.Id)
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "登录成功，欢迎您。",
		"user_id":     u_in_database.Id,
		"token":       t,
	})
}

// 用户信息，我的理解是token是鉴权，看当前是否为登录用户，id才是需要了解的用户详情的id
func UserInfo(c *gin.Context) {
	token := c.Query("token")
	//检查token
	if !utils.CheckToken(token) {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "token已失效",
			"user":        "",
		})
		c.Abort()
		return
	}

	user_id := c.Query("user_id")
	//查询数据库
	u := new(models.User)
	database.Handler.Where("id = ?", user_id).First(u)
	if u.Id == 0 {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  fmt.Sprintf("查询出错，没有id为%v的用户！", user_id),
			"user":        "",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "查询成功",
		"user": gin.H{
			"id":               u.Id,
			"name":             u.Username,
			"follow_count":     0,
			"follower_count":   0,
			"is_follow":        true,
			"avatar":           fmt.Sprintf("%s:%d/public/background/%d.png", config.Server.Host, config.Server.Port, u.Id),
			"background_image": fmt.Sprintf("%s:%d/public/background/%d.png", config.Server.Host, config.Server.Port, u.Id),
			"signature":        "Golang冲冲冲",
			"total_favorited":  "0",
			"work_count":       0,
			"favorite_count":   0,
		},
		//"user": nil,
	})
}

// 投稿
func PublishVideo(c *gin.Context) {
	//获取标题
	title := c.PostForm("title")
	//获取token
	tokenString := c.PostForm("token")
	//检查token
	if !utils.CheckToken(tokenString) {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "token已失效！",
		})
		c.Abort()
		return
	}
	//解析token
	claim, _ := utils.ParseToken(tokenString)
	//检查是否为空
	if title == "" {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "标题不能为空！",
		})
		c.Abort()
		return
	}
	//生成UUID
	u_id := uuid.New()
	key := u_id.String()
	//获取文件
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "视频上传失败！",
		})
		c.Abort()
		return
	}
	//截取文件扩展名
	ext := filepath.Ext(file.Filename)
	//合成文件保存路径：public/uuid.ext
	p := fmt.Sprintf("public/video/%s.%s", key, ext)
	//保存文件
	err = c.SaveUploadedFile(file, "./"+p)
	if err != nil {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "视频保存失败！",
		})
		c.Abort()
		return
	}
	//数据库新建记录
	v := models.Video{
		Path:   p,
		Title:  title,
		UserID: claim.Id,
	}
	database.Handler.Create(&v)
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  fmt.Sprintf("%s视频上传成功。", title),
	})
}
