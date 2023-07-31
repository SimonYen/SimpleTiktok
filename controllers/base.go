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
	"strconv"
	"time"

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
	//生成随机渐变背景图
	utils.GenerateBackground(u.Id)
	//生成GitHub风格头像
	utils.GenerateAvatar(u.Username, u.Id)
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
			"user":        nil,
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
			"user":        nil,
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "查询成功",
		"user": models.UserJSON{
			Avatar:          fmt.Sprintf("%s:%d/public/avatar/%d.png", config.Server.Host, config.Server.Port, u.Id),
			BackgroundImage: fmt.Sprintf("%s:%d/public/background/%d.png", config.Server.Host, config.Server.Port, u.Id),
			FavoriteCount:   0,
			FollowCount:     0,
			FollowerCount:   0,
			ID:              u.Id,
			IsFollow:        false,
			Name:            u.Username,
			Signature:       "Simon冲冲冲",
			TotalFavorited:  "0",
			WorkCount:       0,
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
	//数据库新建记录
	video := models.Video{
		Title:       title,
		UserID:      claim.Id,
		Extension:   ext,
		CreatedTime: time.Now().Unix(),
	}
	database.Handler.Create(&video)
	//合成文件保存路径：public/video/id.ext
	p := fmt.Sprintf("public/video/%d%s", video.Id, ext)
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
	//为视频生成封面
	utils.GenerateCover(video.Id, p)
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  fmt.Sprintf("%s视频上传成功。", title),
	})
}

// 视频流接口
func VideoFeed(c *gin.Context) {
	//感觉token没必要获取
	//获取latest_time
	lt := c.Query("latest_time")
	now := time.Now().Unix()
	latest_time, _ := strconv.Atoi(lt)
	//数据库中查询所有符合条件的video
	var videos []models.Video
	database.Handler.Where("created_time > ? and created_time < ?", latest_time, now).Order("created_time desc").Limit(30).Find(&videos)

	if len(videos) == 0 {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "没有新的视频了",
			"next_time":   nil,
			"video_list":  nil,
		})
		c.Abort()
		return
	}
	video_jsons := make([]models.VideoJSON, 0, 30)
	for _, video := range videos {
		//查询视频作者信息
		var u models.User
		database.Handler.Where("id = ?", video.UserID).First(&u)
		//将相关信息填入结构体中
		v := models.VideoJSON{
			ID: video.Id,
			Author: models.UserJSON{
				Avatar:          fmt.Sprintf("%s:%d/public/avatar/%d.png", config.Server.Host, config.Server.Port, u.Id),
				BackgroundImage: fmt.Sprintf("%s:%d/public/background/%d.png", config.Server.Host, config.Server.Port, u.Id),
				FavoriteCount:   0,
				FollowCount:     0,
				FollowerCount:   0,
				ID:              u.Id,
				IsFollow:        false,
				Name:            u.Username,
				Signature:       "Simon冲冲冲",
				TotalFavorited:  "0",
				WorkCount:       0,
			},
			PlayURL:       fmt.Sprintf("%s:%d/public/video/%d%s", config.Server.Host, config.Server.Port, video.Id, video.Extension),
			CoverURL:      fmt.Sprintf("%s:%d/public/screenshot/%d.png", config.Server.Host, config.Server.Port, video.Id),
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         video.Title,
		}
		video_jsons = append(video_jsons, v)
	}
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "获取视频流成功。",
		"video_list":  video_jsons,
		"next_time":   videos[len(videos)-1].CreatedTime, //要不然不能循环刷
	})
}

func OwnPulishedVideo(c *gin.Context) {
	//获取token
	tokenString := c.Query("token")
	//检查token
	if !utils.CheckToken(tokenString) {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "token已失效！",
			"video_list":  nil,
		})
		c.Abort()
		return
	}
	//解析token
	claim, _ := utils.ParseToken(tokenString)
	//获取user_id
	user_id, _ := strconv.Atoi(c.Query("user_id"))
	if user_id != int(claim.Id) {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "token与本人id不符！",
			"video_list":  nil,
		})
		c.Abort()
		return
	}
	//数据库中查询所有符合条件的video
	var videos []models.Video
	database.Handler.Where("user_id = ?", user_id).Find(&videos)

	if len(videos) == 0 {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "您没有发布视频。",
			"video_list":  nil,
		})
		c.Abort()
		return
	}
	video_jsons := make([]models.VideoJSON, 0, 30)
	for _, video := range videos {
		//查询视频作者信息
		var u models.User
		database.Handler.Where("id = ?", video.UserID).First(&u)
		//将相关信息填入结构体中
		v := models.VideoJSON{
			ID: video.Id,
			Author: models.UserJSON{
				Avatar:          fmt.Sprintf("%s:%d/public/avatar/%d.png", config.Server.Host, config.Server.Port, u.Id),
				BackgroundImage: fmt.Sprintf("%s:%d/public/background/%d.png", config.Server.Host, config.Server.Port, u.Id),
				FavoriteCount:   0,
				FollowCount:     0,
				FollowerCount:   0,
				ID:              u.Id,
				IsFollow:        false,
				Name:            u.Username,
				Signature:       "Simon冲冲冲",
				TotalFavorited:  "0",
				WorkCount:       0,
			},
			PlayURL:       fmt.Sprintf("%s:%d/public/video/%d%s", config.Server.Host, config.Server.Port, video.Id, video.Extension),
			CoverURL:      fmt.Sprintf("%s:%d/public/screenshot/%d.png", config.Server.Host, config.Server.Port, video.Id),
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         video.Title,
		}
		video_jsons = append(video_jsons, v)
	}
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "获取视频流成功。",
		"video_list":  video_jsons,
	})
}
