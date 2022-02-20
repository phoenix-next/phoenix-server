package initialize

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/phoenix-next/phoenix-server/api/v1"
	"github.com/phoenix-next/phoenix-server/middleware"
)

func InitRouter(r *gin.Engine) {
	rawRouter := r.Group("/api/v1")
	rawRouter.Use(cors.Default())

	authRouter := rawRouter.Group("/user")
	{
		authRouter.POST("/register", v1.Register)
		authRouter.POST("/captcha", v1.GetCaptcha)
		authRouter.POST("/login", v1.Login)
	}

	basicRouter := rawRouter.Group("/")
	basicRouter.Use(middleware.AuthRequired())

	userRouter := basicRouter.Group("/user")
	{
		userRouter.GET("/profile", v1.GetProfile)
	}
}
