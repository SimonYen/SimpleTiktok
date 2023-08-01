package routes

import "github.com/gin-gonic/gin"

func Setup(e *gin.Engine) {
	loadBase(e)
	loadInteractive(e)
}
