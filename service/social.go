package service

import (
	"errors"
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
		global.LOG.Panic("GetUserByID: search error")
		return user, true
	} else {
		return user, false
	}
}

// GetUserByEmail 根据用户邮箱查询某个用户
func GetUserByEmail(email string) (user database.User, notFound bool) {
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
func GetAllUser() (users []database.User) {
	users = make([]database.User, 0)
	global.DB.Find(&users)
	return users
}

// CreateCaptcha 生成验证码
func CreateCaptcha(captcha *database.Captcha) (err error) {
	if err = global.DB.Create(captcha).Error; err != nil {
		return err
	}
	return nil
}

// GetCaptchaByEmail 根据邮箱删得到验证码
func GetCaptchaByEmail(name string) (captcha database.Captcha, notFound bool) {
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
	if err = global.DB.Where("email = ?", email).Delete(database.Captcha{}).Error; err != nil {
		return err
	}
	return nil
}

// CreateOrganization 生成组织
func CreateOrganization(organization *database.Organization) (err error) {
	if err = global.DB.Create(organization).Error; err != nil {
		return err
	}
	return nil
}

// DeleteOrganizationByID 根据ID删除组织
func DeleteOrganizationByID(ID uint64) (err error) {
	if err = global.DB.Where("id = ?", ID).Delete(database.Organization{}).Error; err != nil {
		return err
	}
	return nil
}

// GetOrganizationByID 根据组织 ID 查询某个组织
func GetOrganizationByID(ID uint64) (organization database.Organization, notFound bool) {
	err := global.DB.First(&organization, ID).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return organization, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Panic("GetUserByID: search error")
		return organization, true
	} else {
		return organization, false
	}
}

// UpdateOrganization 根据信息更新组织
func UpdateOrganization(organization *database.Organization, name string, profile string) (err error) {
	organization.Name = name
	organization.Profile = profile
	err = global.DB.Save(organization).Error
	return err
}
