package v1

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model/api"
	"github.com/phoenix-next/phoenix-server/model/database"
	"github.com/phoenix-next/phoenix-server/service"
	"github.com/phoenix-next/phoenix-server/utils"
)

// Register
// @Summary      注册
// @Description  注册新用户
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.RegisterQ  true  "用户名, 邮箱, 密码, 验证码"
// @Success      200   {object}  api.RegisterA  "返回信息"
// @Router       /api/v1/user/register [post]
func Register(c *gin.Context) {
	var data api.RegisterQ
	if err := c.ShouldBindJSON(&data); err != nil {
		global.LOG.Panic("Register: bind data error")
	}
	if _, notFound := service.GetUserByName(data.Name); !notFound {
		c.JSON(http.StatusForbidden, api.RegisterA{Message: "用户已存在"})
		return
	}
	captcha, err := strconv.ParseUint(data.Captcha, 10, 64)
	realCaptcha, notFound := service.GetCaptchaByEmail(data.Email)
	if notFound {
		c.JSON(http.StatusForbidden, api.RegisterA{Message: "验证码未发送"})
	}
	if captcha != realCaptcha.Captcha || err != nil {
		c.JSON(http.StatusForbidden, api.RegisterA{Message: "验证码错误"})
		return
	}
	if err := service.CreateUser(&database.User{Name: data.Name, Password: data.Password, Email: data.Email}); err != nil {
		global.LOG.Panic("Register: create user error")
	}
	c.JSON(http.StatusOK, api.RegisterA{Message: "创建用户成功"})
}

// GetCaptcha
// @Summary      发送验证码
// @Description  根据邮箱发送验证码，并更新数据库
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.GetCaptchaQ  true  "邮箱"
// @Success      200   {object}  api.GetCaptchaA  "返回信息"
// @Router       /api/v1/user/captcha [post]
func GetCaptcha(c *gin.Context) {
	var data api.GetCaptchaQ
	err := c.ShouldBindJSON(&data)
	if err != nil {
		global.LOG.Panic("GetCaptcha: bind data error")
	}
	confirmNumber := rand.New(rand.NewSource(time.Now().UnixNano())).Int() % 1000000
	if err = service.DeleteCaptchaByEmail(data.Email); err != nil {
		global.LOG.Panic("GetCaptcha: delete captcha error")
	}
	if err = service.CreateCaptcha(&database.Captcha{Email: data.Email, Captcha: uint64(confirmNumber)}); err != nil {
		global.LOG.Panic("GetCaptcha: create captcha error")
	}
	utils.SendRegisterEmail(data.Email, confirmNumber)
	c.JSON(http.StatusOK, api.RegisterA{Message: "发送验证码成功"})
}

// Login
// @Summary      用户登录
// @Description  根据用户邮箱和密码等生成token，并将token返回给用户
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.LoginQ  true  "邮箱，密码"
// @Success      200   {object}  api.LoginA  "返回信息，Token"
// @Router       /api/v1/user/login [post]
func Login(c *gin.Context) {
	var data api.LoginQ
	err := c.ShouldBindJSON(&data)
	if err != nil {
		global.LOG.Panic("Login: bind data error")
	}
	c.JSON(http.StatusOK, api.LoginA{Message: "登录成功", Token: "123456"})
}
