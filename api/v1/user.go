package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
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
// @Param        data  body      model.RegisterQ  true  "用户名, 邮箱, 密码, 验证码"
// @Success      200   {object}  model.CommonA    "是否成功，返回信息"
// @Router       /api/v1/users [post]
func Register(c *gin.Context) {
	data := utils.BindJsonData(c, &model.RegisterQ{}).(*model.RegisterQ)
	if _, notFound := service.GetUserByEmail(data.Email); !notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "用户已存在"})
		return
	}
	captcha, err := strconv.ParseUint(data.Captcha, 10, 64)
	realCaptcha, notFound := service.GetCaptchaByEmail(data.Email)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "验证码未发送"})
		return
	}
	if captcha != realCaptcha.Captcha || err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "验证码错误"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 12)
	if err != nil {
		global.LOG.Panic("Register: hash password error")
	}
	if err := service.CreateUser(&model.User{Name: data.Name, Password: string(hashedPassword), Email: data.Email}); err != nil {
		global.LOG.Panic("Register: create user error")
	}
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "创建用户成功"})
}

// GetCaptcha
// @Summary      发送验证码
// @Description  根据邮箱发送验证码，并更新数据库
// @Tags         登录模块
// @Accept       json
// @Produce      json
// @Param        data  body      model.GetCaptchaQ  true  "邮箱"
// @Success      200   {object}  model.CommonA      "是否成功，返回信息"
// @Router       /api/v1/captcha [post]
func GetCaptcha(c *gin.Context) {
	data := utils.BindJsonData(c, &model.GetCaptchaQ{}).(*model.GetCaptchaQ)
	confirmNumber := rand.New(rand.NewSource(time.Now().UnixNano())).Int() % 1000000
	if err := service.DeleteCaptchaByEmail(data.Email); err != nil {
		global.LOG.Panic("GetCaptcha: delete captcha error")
	}
	if err := service.CreateCaptcha(&model.Captcha{Email: data.Email, Captcha: uint64(confirmNumber)}); err != nil {
		global.LOG.Panic("GetCaptcha: create captcha error")
	}
	utils.SendRegisterEmail(data.Email, confirmNumber)
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "发送验证码成功"})
}

// Login
// @Summary      用户登录
// @Description  根据用户邮箱和密码等生成token，并将token返回给用户
// @Tags         登录模块
// @Accept       json
// @Produce      json
// @Param        data  body      model.LoginQ  true  "邮箱，密码"
// @Success      200   {object}  model.LoginA  "是否成功，返回信息，Token"
// @Router       /api/v1/tokens [post]
func Login(c *gin.Context) {
	data := utils.BindJsonData(c, &model.LoginQ{}).(*model.LoginQ)
	user, notFound := service.GetUserByEmail(data.Email)
	if notFound {
		c.JSON(http.StatusOK, model.LoginA{Success: false, Message: "登录失败，邮箱不存在", Token: "", ID: 0})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		c.JSON(http.StatusOK, model.LoginA{Success: false, Message: "登录失败，密码错误", Token: "", ID: 0})
		return
	}
	token := utils.GenerateToken(user.ID)
	c.JSON(http.StatusOK, model.LoginA{Success: true, Message: "登录成功", Token: token, ID: user.ID})
}

// GetProfile
// @Summary      获取用户资料
// @Description  获取用户的详细资料
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Success      200      {object}  model.GetProfileA  "是否成功，返回信息，用户资料"
// @Router       /api/v1/users/profile [get]
func GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"TODO": "logic"})
}

// GetUserOrganization
// @Summary      获取用户所属的所有组织
// @Description  获取用户所属的所有组织，返回组织列表
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                      true  "token"
// @Success      200      {object}  model.GetUserOrganizationA  "是否成功，返回信息，用户所属的组织列表"
// @Router       /api/v1/users/organizations [get]
func GetUserOrganization(c *gin.Context) {
	user := utils.SolveUser(c)
	var relation []model.OrganizationT
	global.DB.Model(&model.Invitation{}).Where("user_id = ?", user.ID).Find(&relation)
	c.JSON(http.StatusOK, model.GetUserOrganizationA{Success: true, Message: "获取用户所属组织成功", Organization: relation})
}

// GetUserInvitations
// @Summary      获取用户收到的所有组织邀请
// @Description  组织管理员会邀请用户进入，该接口获得一个用户收到的所有邀请
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                     true  "token"
// @Success      200      {object}  model.GetUserInvitationsA  "是否成功，返回信息，组织信息列表"
// @Router       /api/v1/users/invitations [get]
func GetUserInvitations(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, gin.H{"TODO": "logic"})
}
