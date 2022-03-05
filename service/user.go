package service

import (
	"errors"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"gorm.io/gorm"
)

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

// GetAllUser 查询所有用户
func GetAllUser() (users []model.User) {
	users = make([]model.User, 0)
	global.DB.Find(&users)
	return users
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
