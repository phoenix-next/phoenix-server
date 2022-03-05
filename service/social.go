package service

import (
	"errors"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"gorm.io/gorm"
)

// CreateOrganization 生成组织
func CreateOrganization(organization *model.Organization) (err error) {
	if err = global.DB.Create(organization).Error; err != nil {
		return err
	}
	return nil
}

// DeleteOrganizationByID 根据ID删除组织
func DeleteOrganizationByID(ID uint64) (err error) {
	if err = global.DB.Where("id = ?", ID).Delete(model.Organization{}).Error; err != nil {
		return err
	}
	return nil
}

// GetOrganizationByID 根据组织 ID 查询某个组织
func GetOrganizationByID(ID uint64) (organization model.Organization, notFound bool) {
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

// GetOrganizationByName 根据组织名称查询某个组织
func GetOrganizationByName(name string) (org model.Organization, notFound bool) {
	err := global.DB.Where("name = ?", name).First(&org).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return org, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Panic("GetUserByEmail: search error")
		return org, true
	} else {
		return org, false
	}
}

// UpdateOrganization 根据信息更新组织
func UpdateOrganization(organization *model.Organization, name string, profile string) (err error) {
	organization.Name = name
	organization.Profile = profile
	err = global.DB.Save(organization).Error
	return err
}
