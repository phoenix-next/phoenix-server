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
	"path/filepath"
	"strconv"
	"time"
)

// CreateUser
// @Summary      注册
// @Description  注册新用户
// @Tags         登录模块
// @Accept       json
// @Produce      json
// @Param        data  body      model.CreateUserQ  true  "用户名, 邮箱, 密码, 验证码"
// @Success      200   {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/users [post]
func CreateUser(c *gin.Context) {
	data := utils.BindJsonData(c, &model.CreateUserQ{}).(*model.CreateUserQ)
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
		global.LOG.Panic("CreateUser: hash password error")
	}
	if err := service.CreateUser(&model.User{Name: data.Name, Password: string(hashedPassword), Email: data.Email}); err != nil {
		global.LOG.Panic("CreateUser: create user error")
	}
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "创建用户成功"})
}

// CreateCaptcha
// @Summary      发送验证码
// @Description  根据邮箱发送验证码，并更新数据库
// @Tags         登录模块
// @Accept       json
// @Produce      json
// @Param        data  body      model.CreateCaptchaQ  true  "邮箱"
// @Success      200   {object}  model.CommonA      "是否成功，返回信息"
// @Router       /api/v1/captcha [post]
func CreateCaptcha(c *gin.Context) {
	data := utils.BindJsonData(c, &model.CreateCaptchaQ{}).(*model.CreateCaptchaQ)
	confirmNumber := rand.New(rand.NewSource(time.Now().UnixNano())).Int() % 1000000
	if err := service.DeleteCaptchaByEmail(data.Email); err != nil {
		global.LOG.Panic("CreateCaptcha: delete captcha error")
	}
	if err := service.CreateCaptcha(&model.Captcha{Email: data.Email, Captcha: uint64(confirmNumber)}); err != nil {
		global.LOG.Panic("CreateCaptcha: create captcha error")
	}
	utils.SendRegisterEmail(data.Email, confirmNumber)
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "发送验证码成功"})
}

// CreateToken
// @Summary      用户登录
// @Description  根据用户邮箱和密码等生成token，并将token返回给用户
// @Tags         登录模块
// @Accept       json
// @Produce      json
// @Param        data  body      model.CreateTokenQ  true  "邮箱，密码"
// @Success      200   {object}  model.CreateTokenA  "是否成功，返回信息，Token"
// @Router       /api/v1/tokens [post]
func CreateToken(c *gin.Context) {
	data := utils.BindJsonData(c, &model.CreateTokenQ{}).(*model.CreateTokenQ)
	user, notFound := service.GetUserByEmail(data.Email)
	if notFound {
		c.JSON(http.StatusOK, model.CreateTokenA{Success: false, Message: "登录失败，邮箱不存在", Token: "", ID: 0})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		c.JSON(http.StatusOK, model.CreateTokenA{Success: false, Message: "登录失败，密码错误", Token: "", ID: 0})
		return
	}
	token := utils.GenerateToken(user.ID)
	c.JSON(http.StatusOK, model.CreateTokenA{Success: true, Message: "登录成功", Token: token, ID: user.ID})
}

// UpdateUser
// @Summary      更新用户资料
// @Description  更新发出请求的用户的详细资料
// @Tags         用户模块
// @Accept       multipart/form-data
// @Produce      json
// @Param        x-token  header    string          true  "token"
// @Param        data     body      model.UpdateUserQ  true  "用户名，密码，旧密码，用户简介，用户头像"
// @Success      200      {object}  model.CommonA      "是否成功，返回信息"
// @Router       /api/v1/users [put]
func UpdateUser(c *gin.Context) {
	// 获取请求数据
	user := utils.SolveUser(c)
	// 按照需要更新数据
	if name, found := c.GetPostForm("name"); found {
		user.Name = name
	}
	if oldPassword, found := c.GetPostForm("oldPassword"); found {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
			c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "当前密码输入不正确"})
			return
		}
		password, _ := c.GetPostForm("password")
		if password == "" {
			c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "新密码为空"})
			return
		}
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
		user.Password = string(hashedPassword)
	}
	if profile, found := c.GetPostForm("profile"); found {
		user.Profile = profile
	}
	if avatar, err := c.FormFile("avatar"); err == nil && avatar != nil {
		filename := strconv.FormatUint(user.ID, 10) + "_avatar_" + avatar.Filename
		c.SaveUploadedFile(avatar, filepath.Join(global.VP.GetString("user_path"), filename))
		user.Avatar = "resource/user/" + filename
	}
	// 进行数据库操作并返回
	global.DB.Save(&user)
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新用户信息成功"})
}

// GetUser
// @Summary      获取用户资料
// @Description  获取一个用户公开的详细资料
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       path      int             true  "用户ID"
// @Success      200      {object}  model.GetUserA  "是否成功，返回信息，用户名，用户邮箱，用户头像，用户简介"
// @Router       /api/v1/users/{id} [get]
func GetUser(c *gin.Context) {
	// 获取请求数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.GetUserA{Success: false, Message: "请求参数非法"})
		return
	}
	// 查询对应ID的用户
	user, notFound := service.GetUserByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.GetUserA{Success: false, Message: "找不到对应的用户"})
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, model.GetUserA{
		Success: true,
		Message: "获取用户信息成功",
		Name:    user.Name,
		Email:   user.Email,
		Avatar:  user.Avatar,
		Profile: user.Profile})
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
	global.DB.Model(&model.Invitation{}).Where("user_id = ? and is_valid = ?", user.ID, true).Find(&relation)
	c.JSON(http.StatusOK, model.GetUserOrganizationA{Success: true, Message: "获取用户所属组织成功", Organization: relation})
}

// GetUserInvitation
// @Summary      获取用户收到的所有组织邀请
// @Description  组织管理员会邀请用户进入，该接口用于获取用户未确认的所有邀请
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                    true  "token"
// @Success      200      {object}  model.GetUserInvitationA  "是否成功，返回信息，组织信息列表"
// @Router       /api/v1/users/invitations [get]
func GetUserInvitation(c *gin.Context) {
	user := utils.SolveUser(c)
	var invitations []model.OrganizationT
	global.DB.Model(&model.Invitation{}).Where("user_id = ? AND is_valid = ?", user.ID, false).Find(&invitations)
	c.JSON(http.StatusOK, model.GetUserInvitationA{Success: true, Message: "获取用户未确认邀请成功", Organization: invitations})
}
