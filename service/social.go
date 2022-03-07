package service

import (
	"errors"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"gorm.io/gorm"
)

// Helper

// IsUserInThisOrganization 判断用户是否在该组织中
func IsUserInThisOrganization(uid uint64, orgID uint64) (ok bool, err error) {
	_, notFoundUser := GetUserByID(uid)
	_, notFoundOrg := GetOrganizationByID(orgID)
	if notFoundUser || notFoundOrg {
		return false, errors.New("用户或组织不存在")
	}
	_, notFound := GetInvitationByUserOrg(uid, orgID)
	return !notFound, nil
}

// 数据库操作

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

// CreateInvitation 生成组织邀请
func CreateInvitation(invitation *model.UserOrgRel) (err error) {
	if err = global.DB.Create(invitation).Error; err != nil {
		return err
	}
	return nil
}

// GetInvitationByUserOrg 根据组织与用户查找邀请
func GetInvitationByUserOrg(uid uint64, orgID uint64) (rel *model.UserOrgRel, notFound bool) {
	err := global.DB.Where("uid = ? AND org_id = ? AND IsValid = ?", uid, orgID, true).First(&rel).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return rel, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Panic("GetInvitationByUserOrg: search error")
		return rel, true
	} else {
		return rel, false
	}
}

// UpdateInvitation 根据信息更新邀请
func UpdateInvitation(rel *model.UserOrgRel) (err error) {
	err = global.DB.Save(rel).Error
	return err
}

// GetOrganizationMember 获取组织所有的用户
func GetOrganizationMember(oid uint64) (members []model.Member) {
	var rel *[]model.UserOrgRel
	global.DB.Where("orgID = ? AND IsValid = ?", oid, true).Find(&rel)
	for _, member := range *rel {
		members = append(members, model.Member{ID: member.UserID, Name: member.UserName, IsAdmin: member.IsAdmin})
	}
	return members
}

// GetOrganizationAdmin 获取一个组织中所有管理员的用户ID
func GetOrganizationAdmin(oid uint64) (admin []uint64) {
	global.DB.Model(&model.UserOrgRel{}).Where("orgID = ? AND IsValid = ? AND IsAdmin = ?", oid, true, true).Find(&admin)
	return
}
