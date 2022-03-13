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
		rawRouter.POST("/users", v1.CreateUser)
		rawRouter.POST("/captcha", v1.CreateCaptcha)
		rawRouter.POST("/tokens", v1.CreateToken)
		rawRouter.POST("/password", v1.ResetPassword)
	}
	// 用户头像资源服务器
	avatarRouter := rawRouter.Group("/resource")
	{
		avatarRouter.StaticFS("/user", http.Dir(global.VP.GetString("user_path")))
	}

	// 除了登录模块和头像资源之外，都需要身份认证
	basicRouter := rawRouter.Group("/")
	basicRouter.Use(middleware.AuthRequired())

	// 静态资源服务器
	resourceRouter := basicRouter.Group("/resource")
	{
		resourceRouter.StaticFS("/problem", http.Dir(global.VP.GetString("problem_path")))
		resourceRouter.StaticFS("/tutorial", http.Dir(global.VP.GetString("tutorial_path")))
	}
	// 用户模块
	userRouter := basicRouter.Group("/users")
	{
		userRouter.PUT("", v1.UpdateUser)
		userRouter.GET("/:id", v1.GetUser)
		userRouter.GET("/organizations", v1.GetUserOrganization)
		userRouter.GET("/invitations", v1.GetUserInvitation)
		userRouter.DELETE("/organizations/:id", v1.QuitOrganization)
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
		teamRouter.GET("/:id", v1.GetOrganization)
		teamRouter.DELETE("/:id", v1.DeleteOrganization)
		teamRouter.PUT("/:id", v1.UpdateOrganization)
		teamRouter.POST("/:id/invitations", v1.CreateInvitation)
		teamRouter.POST("/:id/users", v1.UpdateOrganizationMember)
		teamRouter.GET("/:id/users", v1.GetOrganizationMember)
		teamRouter.DELETE("/:id/users/:userID", v1.DeleteOrganizationMember)
		teamRouter.POST("/:id/admins", v1.UpdateOrganizationAdmin)
		teamRouter.DELETE("/:id/admins/:adminID", v1.DeleteOrganizationAdmin)
		teamRouter.GET("/:id/problems", v1.GetOrganizationProblem)
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
	// 教程模块
	tutorialRouter := basicRouter.Group("/tutorials")
	{
		tutorialRouter.GET("", v1.GetTutorialList)
		tutorialRouter.POST("", v1.CreateTutorial)
		tutorialRouter.DELETE("/:id", v1.DeleteTutorial)
		tutorialRouter.GET("/:id", v1.GetTutorial)
		tutorialRouter.PUT("/:id", v1.UpdateTutorial)
		tutorialRouter.GET("/:id/version", v1.GetTutorialVersion)
	}
	// 比赛模块
	contestRouter := basicRouter.Group("/contests")
	{
		contestRouter.GET("", v1.GetContestList)
		contestRouter.POST("", v1.CreateContest)
		contestRouter.GET("/:id", v1.GetContest)
		contestRouter.DELETE("/:id", v1.DeleteContest)
		contestRouter.PUT("/:id", v1.UpdateContest)
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
