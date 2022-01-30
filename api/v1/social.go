package v1

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
// @Success      200   {object}  api.RegisterA  "用户注册账号,返回注册成功信息"
// @Router       /api/v1/user/register [post]
func Register(c *gin.Context) {
	var data api.RegisterQ
	if err := c.ShouldBindJSON(&data); err != nil {
		panic(err)
	}
	if _, notFound := service.GetUserByName(data.Name); !notFound {
		c.JSON(http.StatusForbidden, api.RegisterA{Message: "用户已存在"})
		return
	}
	if !service.ValidEmailCaptcha(data.Email, data.Captcha) {
		c.JSON(http.StatusForbidden, api.RegisterA{Message: "验证码错误"})
		return
	}
	err := service.CreateUser(&database.User{Name: data.Name, Password: data.Password, Email: data.Email})
	if err != nil {
		panic(err)
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
// @Success      200   {object}  api.GetCaptchaA  "用户注册账号,返回注册成功信息"
// @Router       /api/v1/user/captcha [post]
func GetCaptcha(c *gin.Context) {
	var data api.GetCaptchaQ
	if err := c.ShouldBindJSON(&data); err != nil {
		panic(err)
	}
	confirmNumber := rand.New(rand.NewSource(time.Now().UnixNano())).Int() % 1000000
	_ = service.DeleteCaptchaByEmail(data.Email)
	if err := service.CreateCaptcha(&database.Captcha{Email: data.Email, Captcha: confirmNumber}); err != nil {
		panic(err)
	}
	utils.SendRegisterEmail(data.Email, confirmNumber)
	c.JSON(http.StatusOK, api.RegisterA{Message: "发送验证码成功"})
}
