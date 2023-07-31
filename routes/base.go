/*
该文件下保存着基础接口的路由
*/
package routes

import (
	"app/controllers"

	"github.com/gin-gonic/gin"
)

// 注册路由
func loadBase(e *gin.Engine) {
	e.POST("/douyin/user/register/", controllers.Register)
	e.POST("/douyin/user/login/", controllers.Login)
	e.GET("/douyin/user/", controllers.UserInfo)
	e.POST("/douyin/publish/action/", controllers.PublishVideo)
	e.GET("/douyin/feed/", controllers.VideoFeed)
}
