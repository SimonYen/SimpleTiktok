/*
该文件下保存着基础接口的控制器
*/
package controllers

import "github.com/gin-gonic/gin"

func Hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hhh"})
}
