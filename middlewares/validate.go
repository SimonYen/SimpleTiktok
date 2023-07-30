/*
这个包检查query参数是否正常
*/
package middleware

import (
	"app/models"

	"github.com/gin-gonic/gin"
)

func CheckUserQuery(c *gin.Context) {
	u := new(models.User)
	//获取Query参数
	err := c.Bind(u)
	if err != nil {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "解析query参数失败！",
			"user_id":     0,
			"token":       "",
		})
		c.Abort()
		return
	}
	//检查是否为空
	if u.Username == "" || u.Password == "" {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "用户名或者密码不能为空！",
			"user_id":     0,
			"token":       "",
		})
		c.Abort()
		return
	}
	c.Next()
}
