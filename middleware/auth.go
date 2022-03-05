package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/service"
	"github.com/phoenix-next/phoenix-server/utils"
	"net/http"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("x-token")
		id, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户校验失败"})
			c.Abort()
			return
		}
		user, notFound := service.GetUserByID(id)
		if notFound {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户不存在"})
			c.Abort()
			return
		}
		c.Set("user", user)
	}
}
