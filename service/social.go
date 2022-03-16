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

// CreateInvitation 生成组织邀请
func CreateInvitation(invitation *model.Invitation) (err error) {
	if err = global.DB.Create(invitation).Error; err != nil {
		return err
	}
	return nil
}

// GetInValidInvitationByUserOrg 根据组织与用户查找一个待生效邀请
func GetInValidInvitationByUserOrg(uid uint64, orgID uint64) (rel *model.Invitation, notFound bool) {
	err := global.DB.Where("user_id = ? AND org_id = ? AND is_valid = ?", uid, orgID, false).First(&rel).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return rel, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Panic("GetInvitationByUserOrg: search error")
		return rel, true
	} else {
		return rel, false
	}
}

// GetInvitationByUserOrg 根据组织与用户查找一个已生效邀请
func GetInvitationByUserOrg(uid uint64, orgID uint64) (rel *model.Invitation, notFound bool) {
	err := global.DB.Where("user_id = ? AND org_id = ? AND is_valid = ?", uid, orgID, true).First(&rel).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return rel, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Panic("GetInvitationByUserOrg: search error")
		return rel, true
	} else {
		return rel, false
	}
}

// GetOrganizationMember 获取一个组织中所有的用户
func GetOrganizationMember(oid uint64) (members []model.Member) {
	rel := make([]model.Invitation, 0)
	global.DB.Where("org_id = ? AND is_valid = ?", oid, true).Find(&rel)
	for _, member := range rel {
		members = append(members, model.Member{
			ID:      member.UserID,
			Name:    member.UserName,
			Email:   member.UserEmail,
			IsAdmin: member.IsAdmin})
	}
	return members
}

// GetOrganizationAdmin 获取一个组织中所有管理员的Invitation
func GetOrganizationAdmin(oid uint64) (admin []model.Invitation) {
	global.DB.Model(&model.Invitation{}).Where("org_id = ? AND is_valid = ? AND is_admin = ?", oid, true, true).Find(&admin)
	return
}
