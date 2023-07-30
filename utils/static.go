package utils

import (
	"os"

	"github.com/gin-gonic/gin"
)

// 检查文件夹是否存在，存在返回真
func checkIsExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsExist(err)
}

// 加载静态文件
func LoadStatic(e *gin.Engine) {
	//如果目录不存在，就创建它
	if !checkIsExist("./public") {
		os.Mkdir("./public", 0755)
	}
	//设置为静态文件夹
	e.Static("/public", "./public")
}
