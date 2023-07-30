package main

import (
	"app/config"
	"app/routes"
	"app/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//注册路由
	routes.Setup(r)
	//加载静态文件
	utils.LoadStatic(r)

	r.Run(fmt.Sprintf("0.0.0.0:%d", config.Server.Port))
}
