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
	e.POST("/register", controllers.Register)
}
