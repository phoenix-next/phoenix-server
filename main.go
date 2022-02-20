package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/phoenix-next/phoenix-server/docs"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/initialize"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title        PhoeniX API
// @version      1.0
// @description  PhoeniX学习平台的后端API接口

// @contact.name   Matrix53
// @contact.url    https://github.com/matrix53
// @contact.email  1079207272@qq.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	global.VP = initialize.InitViper()
	global.LOG = initialize.InitLogger()
	global.DB = initialize.InitMySQL()

	if !global.VP.Get("server.debug").(bool) {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	initialize.InitRouter(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":" + global.VP.GetString("server.port"))
}
