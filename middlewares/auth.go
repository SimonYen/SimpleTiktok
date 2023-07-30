package middleware

import (
	"app/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Query("token")
		if tokenString == "" {
			ctx.JSON(200, gin.H{
				"status_code": 1,
				"status_msg":  "客户端未携带token",
			})
			ctx.Abort()
			return
		} else {
			claims, err := utils.ParseToken(tokenString)
			if err != nil {
				ctx.JSON(200, gin.H{
					"status_code": 1,
					"status_msg":  "非法token，解析失败！",
				})
				ctx.Abort()
				return
			}
			if claims.ExpiresAt.Time.Unix() < time.Now().Unix() {
				ctx.JSON(200, gin.H{
					"status_code": 1,
					"status_msg":  "登录已过期，请重新登录！",
				})
				ctx.Abort()
				return
			}
			ctx.Next()
		}
	}
}
