package main

import (
	"app/config"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "yzy",
		})
	})
	r.Run(fmt.Sprintf(":%d", config.Conf.ServerPort))
}
