package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"github.com/phoenix-next/phoenix-server/utils"
	"gorm.io/gorm"
	"strconv"
)

// Helper

// JudgeReadPermission 判别用户的可读权限
func JudgeReadPermission(oid uint64, readable int, creatorID uint64, c *gin.Context) bool {
	user := utils.SolveUser(c)
	invitation, notFound := GetInvitationByUserOrg(user.ID, oid)
	switch readable {
	case 0:
		return creatorID == user.ID
	case 1:
		return !notFound && invitation.IsAdmin == true
	case 2:
		return !notFound
	case 3:
		return true
	default:
		return false
	}
}

// JudgeWritePermission 判别用户的可写权限
func JudgeWritePermission(oid uint64, writable int, creatorID uint64, c *gin.Context) bool {
	user := utils.SolveUser(c)
	invitation, notFound := GetInvitationByUserOrg(user.ID, oid)
	switch writable {
	case 0:
		return creatorID == user.ID
	case 1:
		return !notFound && invitation.IsAdmin == true
	default:
		return false
	}
}

// GetAvatarPath 根据用户ID获取用户头像的文件名
func GetAvatarPath(uid uint64) string {
	return strconv.FormatUint(uid, 10) + "_avatar"
}

// 数据库操作

// CreateUser 创建用户
func CreateUser(user *model.User) (err error) {
	if err = global.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByID 根据用户 ID 查询某个用户
func GetUserByID(ID uint64) (user model.User, notFound bool) {
	err := global.DB.First(&user, ID).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Panic("GetUserByID: search error")
		return user, true
	} else {
		return user, false
	}
}

// GetUserByEmail 根据用户邮箱查询某个用户
func GetUserByEmail(email string) (user model.User, notFound bool) {
	err := global.DB.Where("email = ?", email).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Panic("GetUserByEmail: search error")
		return user, true
	} else {
		return user, false
	}
}

// CreateCaptcha 生成验证码
func CreateCaptcha(captcha *model.Captcha) (err error) {
	if err = global.DB.Create(captcha).Error; err != nil {
		return err
	}
	return nil
}

// GetCaptchaByEmail 根据邮箱删得到验证码
func GetCaptchaByEmail(name string) (captcha model.Captcha, notFound bool) {
	err := global.DB.Where("email = ?", name).First(&captcha).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return captcha, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Panic("GetCaptchaByEmail: search error")
		return captcha, true
	} else {
		return captcha, false
	}
}

// DeleteCaptchaByEmail 根据邮箱删除验证码
func DeleteCaptchaByEmail(email string) (err error) {
	if err = global.DB.Where("email = ?", email).Delete(model.Captcha{}).Error; err != nil {
		return err
	}
	return nil
}

// GetAdminOrganization 获取用户所有已生效的Invitation，且在这些Invitation中该用户为管理员
func GetAdminOrganization(uid uint64) (admin []model.Invitation) {
	global.DB.Model(&model.Invitation{}).Where("user_id = ? AND is_valid = ? AND is_admin = ?", uid, true, true).Find(&admin)
	return
}

// GetUserOrganization 获取一个用户的所有已生效Invitation
func GetUserOrganization(uid uint64) (invitations []model.Invitation) {
	global.DB.Model(&model.Invitation{}).Where("user_id = ? AND is_valid = ?", uid, true).Find(&invitations)
	return
}

// GetUserInvitation 获取一个用户的所有未生效Invitation
func GetUserInvitation(uid uint64) (invitations []model.Invitation) {
	global.DB.Model(&model.Invitation{}).Where("user_id = ? AND is_valid = ?", uid, false).Find(&invitations)
	return
}
