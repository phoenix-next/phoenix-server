package v1

import (
	"crypto/md5"
	"fmt"
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
	// 获取请求数据
	data := utils.BindJsonData(c, &model.CreateUserQ{}).(*model.CreateUserQ)
	// 用户的邮箱已经注册过的情况
	if _, notFound := service.GetUserByEmail(data.Email); !notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "用户已存在"})
		return
	}
	// 验证码还未发送的情况
	var realCaptcha model.Captcha
	if global.DB.Where("email = ? AND type = ?", data.Email, 1).Find(&realCaptcha).Error != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "验证码未发送"})
		return
	}
	// 验证码不正确的情况
	captcha, err := strconv.ParseUint(data.Captcha, 10, 64)
	if captcha != realCaptcha.Captcha || err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "验证码错误"})
		return
	}
	// 将密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 12)
	if err != nil {
		global.LOG.Panic("CreateUser: hash password error")
	}
	// 成功创建用户
	if err := service.CreateUser(&model.User{Name: data.Name, Password: string(hashedPassword), Email: data.Email}); err != nil {
		global.LOG.Panic("CreateUser: create user error")
	}
	// 删除验证码，防止用户利用
	global.DB.Delete(&realCaptcha)
	// 返回响应
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "注册成功"})
}

// CreateCaptcha
// @Summary      发送验证码
// @Description  根据邮箱发送验证码，类型1为注册验证码，类型2为忘记密码验证码
// @Tags         登录模块
// @Accept       json
// @Produce      json
// @Param        data  body      model.CreateCaptchaQ  true  "邮箱，验证码类型"
// @Success      200   {object}  model.CommonA      "是否成功，返回信息"
// @Router       /api/v1/captcha [post]
func CreateCaptcha(c *gin.Context) {
	// 获取请求数据
	data := utils.BindJsonData(c, &model.CreateCaptchaQ{}).(*model.CreateCaptchaQ)
	// 生成验证码
	confirmNumber := rand.New(rand.NewSource(time.Now().UnixNano())).Int() % 1000000
	// 删除旧的验证码
	global.DB.Where("type = ? AND email = ?", data.Type, data.Email).Delete(&model.Captcha{})
	// 生成新的验证码
	captcha := model.Captcha{Email: data.Email, Captcha: uint64(confirmNumber), Type: data.Type}
	if err := service.CreateCaptcha(&captcha); err != nil {
		global.LOG.Panic("CreateCaptcha: create captcha error")
	}
	// 按照类型发送验证码
	if data.Type == 1 {
		utils.SendRegisterEmail(data.Email, confirmNumber)
	} else {
		utils.SendResetEmail(data.Email, confirmNumber)
	}
	// 返回响应
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "成功发送验证码"})
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
	// 获取请求中的数据
	data := utils.BindJsonData(c, &model.CreateTokenQ{}).(*model.CreateTokenQ)
	// 用于登录的邮箱未注册的情况
	user, notFound := service.GetUserByEmail(data.Email)
	if notFound {
		c.JSON(http.StatusOK, model.CreateTokenA{Success: false, Message: "登录失败，邮箱不存在", Token: "", ID: 0})
		return
	}
	// 密码错误的情况
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		c.JSON(http.StatusOK, model.CreateTokenA{Success: false, Message: "登录失败，密码错误", Token: "", ID: 0})
		return
	}
	// 成功返回响应
	token := utils.GenerateToken(user.ID)
	c.JSON(http.StatusOK, model.CreateTokenA{Success: true, Message: "登录成功", Token: token, ID: user.ID})
}

// UpdateUser
// @Summary      更新用户资料
// @Description  更新发出请求的用户的详细资料
// @Tags         用户模块
// @Accept       multipart/form-data
// @Produce      json
// @Param        x-token  header    string             true   "token"
// @Param        avatar   formData  file               false  "用户头像"
// @Param        data     body      model.UpdateUserQ  false  "用户名，密码，旧密码，用户简介"
// @Success      200   {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/users [put]
func UpdateUser(c *gin.Context) {
	// 获取请求数据
	user := utils.SolveUser(c)
	// 更新用户名
	if name, found := c.GetPostForm("name"); found {
		user.Name = name
		// 维护成员 - 组织关系
		global.DB.Model(&model.Invitation{}).Where("user_id = ?", user.ID).Update("user_name", name)
	}
	// 更新用户密码
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
	// 更新用户简介
	if profile, found := c.GetPostForm("profile"); found {
		user.Profile = profile
	}
	// 更新用户头像
	if avatar, err := c.FormFile("avatar"); err == nil && avatar != nil {
		filename := "user_" + strconv.FormatUint(user.ID, 10) + "_avatar_" + avatar.Filename
		_ = c.SaveUploadedFile(avatar, filepath.Join(global.VP.GetString("image_path"), filename))
		user.Avatar = "resource/image/" + filename
	}
	// 进行数据库操作并返回
	global.DB.Save(&user)
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新成功"})
}

// GetUser
// @Summary      获取用户资料
// @Description  获取一个用户公开的详细资料
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string          true  "token"
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
// @Param        x-token  header    string         true  "token"
// @Success      200      {object}  model.GetUserOrganizationA  "是否成功，返回信息，用户所属的组织列表"
// @Router       /api/v1/users/organizations [get]
func GetUserOrganization(c *gin.Context) {
	// 获取请求中的数据
	user := utils.SolveUser(c)
	// 查询数据库
	relations := make([]model.Invitation, 0)
	global.DB.Where("user_id = ? and is_valid = ?", user.ID, true).Find(&relations)
	// 关系中存入组织头像
	finalRelations := make([]model.OrganizationT, 0)
	for _, rel := range relations {
		org, _ := service.GetOrganizationByID(rel.OrgID)
		finalRelations = append(finalRelations, model.OrganizationT{
			IsAdmin: rel.IsAdmin,
			OrgID:   rel.OrgID,
			Avatar:  org.Avatar,
			OrgName: org.Name})
	}
	// 返回响应
	c.JSON(http.StatusOK, model.GetUserOrganizationA{Success: true, Organization: finalRelations})
}

// QuitOrganization
// @Summary      用户主动退出组织
// @Description  用户主动退出一个已加入的组织
// @Tags         用户模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                      true  "token"
// @Param        id       path      int            true  "组织ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/users/organizations/{id} [delete]
func QuitOrganization(c *gin.Context) {
	// 获取请求参数
	user := utils.SolveUser(c)
	oid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "请求参数非法"})
		return
	}
	// 成员不属于该组织的情况
	rel, notFound := service.GetInvitationByUserOrg(user.ID, oid)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "未加入该组织"})
		return
	}
	// 组织创建者不能退出该组织,只能解散
	org, _ := service.GetOrganizationByID(oid)
	if org.CreatorID == user.ID {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "组织创建者不能退出组织"})
		return
	}
	// 退出组织成功
	global.DB.Delete(&rel)
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "已退出组织"})
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
	// 获取请求中的数据
	user := utils.SolveUser(c)
	// 查询数据库
	invitations := make([]model.Invitation, 0)
	global.DB.Where("user_id = ? AND is_valid = ?", user.ID, false).Find(&invitations)
	// 关系中存入组织头像
	finalInvitations := make([]model.OrganizationT, 0)
	for _, rel := range invitations {
		org, _ := service.GetOrganizationByID(rel.OrgID)
		finalInvitations = append(finalInvitations, model.OrganizationT{
			IsAdmin: rel.IsAdmin,
			OrgID:   rel.OrgID,
			Avatar:  org.Avatar,
			OrgName: org.Name})
	}
	// 返回响应
	c.JSON(http.StatusOK, model.GetUserInvitationA{Success: true, Organization: finalInvitations})
}

// ResetPassword
// @Summary      忘记密码
// @Description  用户忘记密码，根据邮箱验证码重新设置密码
// @Tags         登录模块
// @Accept       json
// @Produce      json
// @Param        data  body      model.ResetPasswordQ  true  "用户邮箱，新密码，验证码"
// @Success      200      {object}  model.CommonA      "是否成功，返回信息"
// @Router       /api/v1/password [post]
func ResetPassword(c *gin.Context) {
	// 获取请求中的数据
	data := utils.BindJsonData(c, &model.ResetPasswordQ{}).(*model.ResetPasswordQ)
	// 找不到验证码的情况
	var realCaptcha model.Captcha
	err := global.DB.Where("email = ? AND type = ?", data.Email, 2).Find(&realCaptcha).Error
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "验证码未发送"})
		return
	}
	// 验证码错误的情况
	captcha, err := strconv.ParseUint(data.Captcha, 10, 64)
	if captcha != realCaptcha.Captcha || err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "验证码错误"})
		return
	}
	// 对新密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 12)
	if err != nil {
		global.LOG.Panic("ResetPassword: hash password error")
	}
	// 重置用户的密码
	user, _ := service.GetUserByEmail(data.Email)
	user.Password = string(hashedPassword)
	if err = global.DB.Save(user).Error; err != nil {
		global.LOG.Panic("ResetPassword: save password error")
	}
	// 删除验证码，防止用户利用验证码
	global.DB.Delete(&realCaptcha)
	// 返回响应
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "已重置密码"})
}

// UploadImage
// @Summary      上传图片
// @Description  用户上传一张图片，用于展示在题目或教程中
// @Tags         用户模块
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file                true  "图片文件"
// @Success      200    {object}  model.UploadImageA  "是否成功，返回信息，生成的图片路径"
// @Router       /api/v1/resource/image [post]
func UploadImage(c *gin.Context) {
	// 获取请求数据
	user := utils.SolveUser(c)
	image, err := c.FormFile("image")
	// 上传失败的情况
	if err != nil || image == nil {
		c.JSON(http.StatusOK, model.UploadImageA{Success: false, Message: "上传图片失败"})
		return
	}
	// md5 哈希
	checksum := md5.Sum([]byte(time.Now().Format("20060102150405") + strconv.FormatUint(user.ID, 10) + "_" + image.Filename))
	filename := fmt.Sprintf("%x", checksum) + filepath.Ext(image.Filename)
	_ = c.SaveUploadedFile(image, filepath.Join(global.VP.GetString("image_path"), filename))
	// 返回文件路径
	c.JSON(http.StatusOK, model.UploadImageA{Success: true, Message: "上传成功", ImagePath: "resource/image/" + filename})
}
