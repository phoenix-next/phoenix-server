package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/utils"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("x-token")
		email, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户校验失败"})
			c.Abort()
			return
		}
		c.Set("email", email)
	}
}
