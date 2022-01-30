package initialize

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/phoenix-next/phoenix-server/api/v1"
)

func InitRouter(r *gin.Engine) {
	basicRouter := r.Group("/api/v1")

	userRouter := basicRouter.Group("/user")
	{
		userRouter.POST("/register", v1.Register)
		userRouter.POST("/captcha", v1.CaptchaValid)
	}
}
