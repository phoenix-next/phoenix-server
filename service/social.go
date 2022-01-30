package service

import (
	"errors"
	"fmt"
	"github.com/phoenix-next/phoenix-server/middleware"
	"strconv"

	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model/database"
	"gorm.io/gorm"
)

// CreateUser 创建用户
func CreateUser(user *database.User) (err error) {
	if err = global.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByID 根据用户 ID 查询某个用户
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

// GetUserByName 根据用户 name 查询某个用户
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

// GetAllUser 查询所有用户
func GetAllUser() (users []database.User) {
	users = make([]database.User, 0)
	global.DB.Find(&users)
	return users
}

// CreateCaptcha CreateUser 创建用户
func CreateCaptcha(captcha *database.Captcha) (err error) {
	if err = global.DB.Create(captcha).Error; err != nil {
		return err
	}
	return nil
}

// GetCaptchaByNEmail 根据邮箱删得到验证码
func GetCaptchaByNEmail(name string) (captcha database.Captcha, notFound bool) {
	err := global.DB.Where("email = ?", name).First(&captcha).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return captcha, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return captcha, false
	}
}

// DeleteCaptchaByEmail 根据邮箱删除验证码
func DeleteCaptchaByEmail(email string) (err error) {
	if err = global.DB.Where("email = ?", email).Delete(database.Captcha{}).Error; err != nil {
		return err
	}
	return nil
}

// SendRegisterEmail 发送注册验证码邮件
func SendRegisterEmail(themail string, number int) {
	subject := "欢迎注册phoenix    xxxx代填"
	// 邮件正文
	mailTo := []string{
		themail,
	}
	body := "Hello,This is a email,这是你的注册码" + strconv.Itoa(number)
	err := middleware.SendMail(mailTo, subject, body)
	if err != nil {
		panic(err)
	}
	fmt.Println("sendRegisterEmail successfully")
	return
}

// 验证邮箱验证码是否正确
func ValidEmailCaptcha(email string, number int) bool {
	if captcha, notFound := GetCaptchaByNEmail(email); !notFound {
		return captcha.Captcha == number
	}
	return false
}
