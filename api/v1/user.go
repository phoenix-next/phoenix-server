package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model/api"
	"github.com/phoenix-next/phoenix-server/model/database"
	"github.com/phoenix-next/phoenix-server/service"
	"github.com/phoenix-next/phoenix-server/utils"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Register
// @Summary      注册
// @Description  注册新用户
// @Tags         登录模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.RegisterQ  true  "用户名, 邮箱, 密码, 验证码"
// @Success      200   {object}  api.CommonA    "是否成功，返回信息"
// @Router       /api/v1/users [post]
func Register(c *gin.Context) {
	var data api.RegisterQ
	if err := c.ShouldBindJSON(&data); err != nil {
		global.LOG.Panic("Register: bind data error")
	}
	if _, notFound := service.GetUserByEmail(data.Email); !notFound {
		c.JSON(http.StatusOK, api.CommonA{Success: false, Message: "用户已存在"})
		return
	}
	captcha, err := strconv.ParseUint(data.Captcha, 10, 64)
	realCaptcha, notFound := service.GetCaptchaByEmail(data.Email)
	if notFound {
		c.JSON(http.StatusOK, api.CommonA{Success: false, Message: "验证码未发送"})
		return
	}
	if captcha != realCaptcha.Captcha || err != nil {
		c.JSON(http.StatusOK, api.CommonA{Success: false, Message: "验证码错误"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 12)
	if err != nil {
		global.LOG.Panic("Register: hash password error")
	}
	if err := service.CreateUser(&database.User{Name: data.Name, Password: string(hashedPassword), Email: data.Email}); err != nil {
		global.LOG.Panic("Register: create user error")
	}
	c.JSON(http.StatusOK, api.CommonA{Success: true, Message: "创建用户成功"})
}

// GetCaptcha
// @Summary      发送验证码
// @Description  根据邮箱发送验证码，并更新数据库
// @Tags         登录模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.GetCaptchaQ  true  "邮箱"
// @Success      200   {object}  api.CommonA      "是否成功，返回信息"
// @Router       /api/v1/captcha [post]
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
	c.JSON(http.StatusOK, api.CommonA{Success: true, Message: "发送验证码成功"})
}

// Login
// @Summary      用户登录
// @Description  根据用户邮箱和密码等生成token，并将token返回给用户
// @Tags         登录模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.LoginQ  true  "邮箱，密码"
// @Success      200   {object}  api.LoginA  "是否成功，返回信息，Token"
// @Router       /api/v1/tokens [post]
func Login(c *gin.Context) {
	var data api.LoginQ
	err := c.ShouldBindJSON(&data)
	if err != nil {
		global.LOG.Panic("Login: bind data error")
	}
	user, notFound := service.GetUserByEmail(data.Email)
	if notFound {
		c.JSON(http.StatusOK, api.LoginA{Success: false, Message: "登录失败，邮箱不存在", Token: "", ID: 0})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		c.JSON(http.StatusOK, api.LoginA{Success: false, Message: "登录失败，密码错误", Token: "", ID: 0})
		return
	}
	token := utils.GenerateToken(user.Email)
	c.JSON(http.StatusOK, api.LoginA{Success: true, Message: "登录成功", Token: token, ID: user.ID})
}

// GetProfile
// @Summary      获取用户资料
// @Description  根据用户ID，获取用户资料
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       path      string           true  "用户ID"
// @Success      200      {object}  api.GetProfileA  "是否成功，返回信息，用户资料"
// @Router       /api/v1/users/{id}/profile [get]
func GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, c.GetString("email"))
}

// GetUserOrganization
// @Summary      获取用户所属的所有组织
// @Description  根据一个用户的ID，获取用户所属的所有组织
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                    true  "token"
// @Param        id       path      int                       true  "用户ID"
// @Success      200      {object}  api.GetUserOrganizationA  "是否成功，返回信息，用户所属的组织列表"
// @Router       /api/v1/users/{id}/organizations [get]
func GetUserOrganization(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// GetUserInvitations
// @Summary      获取用户收到的所有组织邀请
// @Description  组织管理员会邀请用户进入，该接口获得一个用户收到的所有邀请
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                   true  "token"
// @Param        id       path      int                true  "用户ID"
// @Success      200      {object}  api.GetUserInvitationsA  "是否成功，返回信息，组织信息列表"
// @Router       /api/v1/users/{id}/invitations [get]
func GetUserInvitations(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// GetAdminInfo
// @Summary      获取用户的所有管理员权限信息
// @Description  用户是哪些组织的管理员？返回这些组织的基础信息
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string           true  "token"
// @Param        id       path      int                      true  "用户ID"
// @Success      200      {object}  api.GetAdminInfoA  "是否成功，返回信息，组织信息列表"
// @Router       /api/v1/users/{id}/admins [get]
func GetAdminInfo(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}
