/*
该文件下保存着基础接口的路由
*/
package routes

import (
	"app/controllers"
	middleware "app/middlewares"

	"github.com/gin-gonic/gin"
)

// 注册路由
func loadBase(e *gin.Engine) {
	e.POST("/douyin/user/register/", middleware.CheckUserQuery, controllers.Register)
	e.POST("/douyin/user/login/", middleware.CheckUserQuery, controllers.Login)
	e.GET("/douyin/user/", middleware.JWT(), controllers.UserInfo)
}
