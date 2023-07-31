package utils

import (
	"os"

	"github.com/gin-gonic/gin"
)

// 检查文件夹是否存在，不存在就创建它
func createdIfNotExist(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
}

// 加载静态文件
func LoadStatic(e *gin.Engine) {
	//如果目录不存在，就创建它
	createdIfNotExist("./public")
	createdIfNotExist("./public/avatar")
	createdIfNotExist("./public/background")
	createdIfNotExist("./public/video")
	//设置为静态文件夹
	e.Static("/public", "./public")
}
