package v1

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/model/api"
	"github.com/phoenix-next/phoenix-server/model/database"
	"github.com/phoenix-next/phoenix-server/service"
)

// Register      注册
// @Description  注册新用户
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.RegisterQ  true  "用户名, 邮箱, 密码, 验证码"
// @Success      200   {object}  api.RegisterA  "用户注册账号,返回注册成功信息"
// @Router       /api/v1/user/register [post]
func Register(c *gin.Context) {
	var data api.RegisterQ
	if err := c.ShouldBindJSON(&data); err != nil {
		panic(err)
	}
	if _, notFound := service.GetUserByName(data.Name); !notFound {
		c.JSON(http.StatusOK, api.RegisterA{Message: "用户已存在", Code: 401})
		return
	}
	if !service.ValidEmailCaptcha(data.Email, data.Captcha) {
		c.JSON(http.StatusOK, api.RegisterA{Message: "验证码错误", Code: 402})
		return
	}
	err := service.CreateUser(&database.User{Name: data.Name, Password: data.Password, Email: data.Email})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, api.RegisterA{Message: "创建用户成功", Code: 200})
}

// CaptchaValid      邮箱验证
// @Description  根据邮箱发送验证吗，并更新数据库
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.CaptchaValidQ  true  "邮箱"
// @Success      200   {object}  api.RegisterA  "用户注册账号,返回注册成功信息"
// @Router       /api/v1/user/captcha [post]
func CaptchaValid(c *gin.Context) {
	var data api.CaptchaValidQ
	if err := c.ShouldBindJSON(&data); err != nil {
		panic(err)
	}
	confirmNumber := rand.New(rand.NewSource(time.Now().UnixNano())).Int() % 1000000
	_ = service.DeleteCaptchaByEmail(data.Email)
	if err := service.CreateCaptcha(&database.Captcha{Email: data.Email, Captcha: confirmNumber}); err != nil {
		panic(err)
	}
	service.SendRegisterEmail(data.Email, confirmNumber)
	c.JSON(http.StatusOK, api.RegisterA{Message: "发送验证码成功", Code: 200})
}
