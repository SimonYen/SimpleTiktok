/*
该文件下保存着互动接口的路由
*/
package routes

import (
	"app/controllers"

	"github.com/gin-gonic/gin"
)

func loadInteractive(e *gin.Engine) {
	e.POST("/douyin/favorite/action/", controllers.VideoLike)
}
