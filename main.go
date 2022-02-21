package main

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/initialize"
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
	// 初始化全局资源
	global.VP = initialize.InitViper()
	global.LOG = initialize.InitLogger()
	global.DB = initialize.InitMySQL()
	// 创建Router
	if !global.VP.GetBool("server.debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	// 初始化Router
	initialize.InitRouter(r)
	// 运行Router
	initialize.RunRouter(r)
}
