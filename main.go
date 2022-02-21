package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/phoenix-next/phoenix-server/docs"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/initialize"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"os"
	"path/filepath"
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
	isDebug := global.VP.Get("server.debug").(bool)
	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	// 初始化Router
	err := r.SetTrustedProxies(nil)
	if err != nil {
		global.LOG.Panic("初始化失败：禁止使用代理访问失败")
	}
	initialize.InitRouter(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 运行Router
	if isDebug {
		err = r.Run(":" + global.VP.GetString("server.port"))
	} else {
		var path string
		path, err = os.Executable()
		if err != nil {
			global.LOG.Panic("初始化失败：可执行程序路径获取失败")
		}
		folder := filepath.Dir(path)
		err = r.RunTLS(":"+global.VP.GetString("server.port"),
			filepath.Join(folder, global.VP.GetString("server.cert")),
			filepath.Join(folder, global.VP.GetString("server.key")))
	}
	// Router运行错误处理
	if err != nil {
		global.LOG.Panic("运行时错误：", err)
	}
}
