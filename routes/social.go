/*
该文件下保存着社交接口的路由
*/
package routes

import (
	"app/controllers"

	"github.com/gin-gonic/gin"
)

func loadSocial(e *gin.Engine) {
	e.POST("/douyin/relation/action/", controllers.Follow)
	e.GET("/douyin/relation/follow/list/", controllers.GetFollowList)
	e.GET("/douyin/relation/follower/list/", controllers.GetFansList)
	e.GET("/douyin/relation/friend/list/", controllers.GetFridendList)
}
