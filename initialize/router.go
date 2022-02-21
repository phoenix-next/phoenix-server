package initialize

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/phoenix-next/phoenix-server/api/v1"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/middleware"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/unrolled/secure"
	"net/http"
	"os"
	"path/filepath"
)

func InitRouter(r *gin.Engine) {
	r.Use(cors.Default())

	if err := r.SetTrustedProxies(nil); err != nil {
		global.LOG.Panic("初始化失败：禁止使用代理访问失败")
	}

	if global.VP.GetBool("server.docs") {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	rawRouter := r.Group("/api/v1")
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

func RunRouter(r *gin.Engine) {
	// 是否以Debug模式运行
	if global.VP.GetBool("server.debug") {
		global.LOG.Panic("运行时错误：", r.Run(":"+global.VP.GetString("server.port")))
	} else {
		// 将API文档页面重定向至https
		if global.VP.GetBool("server.docs") {
			// 新建请求处理器用于重定向
			redirect := secure.New(secure.Options{
				SSLRedirect: true,
			}).Handler(nil)
			// 创建go线程监听http请求，并用新建的请求处理器重定向
			go global.LOG.Panic("运行时错误：", http.ListenAndServe(":80", redirect))
		}
		// 获取配置并运行Router
		var path string
		path, err := os.Executable()
		if err != nil {
			global.LOG.Panic("初始化失败：可执行程序路径获取失败")
		}
		folder := filepath.Dir(path)
		global.LOG.Panic("运行时错误：",
			r.RunTLS(":"+global.VP.GetString("server.port"),
				filepath.Join(folder, global.VP.GetString("server.cert")),
				filepath.Join(folder, global.VP.GetString("server.key"))))
	}
}
