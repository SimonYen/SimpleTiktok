package main

import (
	"app/config"
	"app/routes"
	"app/utils"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	//用于随机渐变背景生成，只需要调用一次
	rand.Seed(time.Now().Unix())

	r := gin.Default()
	//注册路由
	routes.Setup(r)
	//加载静态文件
	utils.LoadStatic(r)

	r.Run(fmt.Sprintf("0.0.0.0:%d", config.Server.Port))
}
