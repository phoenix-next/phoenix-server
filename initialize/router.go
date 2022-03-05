package initialize

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/phoenix-next/phoenix-server/api/v1"
	_ "github.com/phoenix-next/phoenix-server/docs"
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
	// 跨域配置
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "x-token")
	r.Use(cors.New(config))
	// 是否开启api文档页面
	if global.VP.GetBool("server.docs") {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	// 禁用代理访问
	if err := r.SetTrustedProxies(nil); err != nil {
		global.LOG.Panic("初始化失败：禁止使用代理访问失败")
	}

	// 登录模块
	rawRouter := r.Group("/api/v1")
	{
		rawRouter.POST("/users", v1.Register)
		rawRouter.POST("/captcha", v1.GetCaptcha)
		rawRouter.POST("/tokens", v1.Login)
	}

	// 上传文件模块 查看题目不需要登录即查看
	resourceRouter := rawRouter.Group("/resource")
	{
		resourceRouter.StaticFS("/problem", http.Dir(filepath.Join(global.VP.GetString("root_path"), "resource", "problems")))
	}

	// 除了登录模块之外，都需要身份认证
	basicRouter := rawRouter.Group("/")
	basicRouter.Use(middleware.AuthRequired())

	// 用户模块
	userRouter := basicRouter.Group("/users")
	{
		userRouter.GET("/:id/profile", v1.GetProfile)
		userRouter.GET("/:id/organizations", v1.GetUserOrganization)
		userRouter.GET("/:id/invitations", v1.GetUserInvitations)
		userRouter.GET("/:id/admins", v1.GetAdminInfo)
	}

	// 评测模块
	problemRouter := basicRouter.Group("/problems")
	{
		problemRouter.GET("", v1.GetProblemList)
		problemRouter.POST("", v1.CreateProblem)
		problemRouter.DELETE("/:id", v1.DeleteProblem)
		problemRouter.GET("/:id", v1.GetProblem)
		problemRouter.PUT("/:id", v1.UpdateProblem)
		problemRouter.GET("/:id/version", v1.GetProblemVersion)
	}

	// 组织模块
	teamRouter := basicRouter.Group("/organizations")
	{
		teamRouter.POST("", v1.CreateOrganization)
		teamRouter.DELETE("/:id", v1.DeleteOrganization)
		teamRouter.PUT("/:id", v1.UpdateOrganization)
		teamRouter.POST("/:id/invitations", v1.CreateInvitation)
		teamRouter.POST("/:id/users", v1.UpdateOrganizationMember)
		teamRouter.GET("/:id/users", v1.GetOrganizationMember)
		teamRouter.POST("/:id/admins", v1.UpdateOrganizationAdmin)
		teamRouter.DELETE("/:id/admins/:adminID", v1.DeleteOrganizationAdmin)
	}

	// 论坛模块
	forumRouter := basicRouter.Group("")
	{
		forumRouter.POST("/posts", v1.CreatePost)
		forumRouter.DELETE("/posts/:id", v1.DeletePost)
		forumRouter.PUT("/posts/:id", v1.UpdatePost)
		forumRouter.GET("/posts/:id", v1.GetPost)
		forumRouter.GET("/posts", v1.GetAllPost)
		forumRouter.POST("/posts/:id/comments", v1.CreateComment)
		forumRouter.PUT("/comments/:id", v1.UpdateComment)
		forumRouter.DELETE("/comments/:id", v1.DeleteComment)
		forumRouter.GET("posts/:id/comments", v1.GetComment)
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
			go func() { global.LOG.Panic("运行时错误：", http.ListenAndServe(":80", redirect)) }()
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
