package service

import (
	"errors"

	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model/database"
	"gorm.io/gorm"
)

// 创建用户
func CreateUser(user *database.User) (err error) {
	if err = global.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// 根据用户 ID 查询某个用户
func GetUserByID(ID uint64) (user database.User, notFound bool) {
	err := global.DB.First(&user, ID).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return user, false
	}
}

// 根据用户 name 查询某个用户
func GetUserByName(name string) (user database.User, notFound bool) {
	err := global.DB.Where("name = ?", name).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return user, false
	}
}

// 查询所有用户
func GetAllUser() (users []database.User) {
	users = make([]database.User, 0)
	global.DB.Find(&users)
	return users
}

// 创建用户
func CreateCaptcha(captcha *database.Captcha) (err error) {
	if err = global.DB.Create(captcha).Error; err != nil {
		return err
	}
	return nil
}

// 根据邮箱删得到验证码
func GetCaptchaByEmail(name string) (captcha database.Captcha, notFound bool) {
	err := global.DB.Where("email = ?", name).First(&captcha).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return captcha, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return captcha, false
	}
}

// 根据邮箱删除验证码
func DeleteCaptchaByEmail(email string) (err error) {
	if err = global.DB.Where("email = ?", email).Delete(database.Captcha{}).Error; err != nil {
		return err
	}
	return nil
}

// 验证邮箱验证码是否正确
func ValidEmailCaptcha(email string, number int) bool {
	if captcha, notFound := GetCaptchaByEmail(email); !notFound {
		return captcha.Captcha == number
	}
	return false
}
