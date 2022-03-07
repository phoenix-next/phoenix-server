package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"github.com/phoenix-next/phoenix-server/service"
	"github.com/phoenix-next/phoenix-server/utils"
	"net/http"
	"strconv"
)

// CreateOrganization
// @Summary      创建一个组织
// @Description  创建一个组织，创建者比管理员权限高，也算是管理员
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                     true  "token"
// @Param        data     body      model.CreateOrganizationQ  true  "组织名称，组织的简介"
// @Success      200      {object}  model.CommonA              "是否成功，返回信息"
// @Router       /api/v1/organizations [post]
func CreateOrganization(c *gin.Context) {
	var data model.CreateOrganizationQ
	if err := c.ShouldBindJSON(&data); err != nil {
		global.LOG.Panic("CreateOrganization: bind data error")
	}
	if _, notFound := service.GetOrganizationByName(data.Name); !notFound {
		global.LOG.Warn("CreateOrganization: find same organization name")
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "已存在该名称的组织"})
		return
	}
	user := utils.SolveUser(c)
	organization := model.Organization{Name: data.Name, Profile: data.Profile, CreatorName: user.Name, CreatorID: user.ID}
	if err := service.CreateOrganization(&organization); err != nil {
		global.LOG.Panic("CreateOrganization: create organization error")
	}
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "创建组织成功"})
}

// UpdateOrganization
// @Summary      更新一个组织的信息
// @Description  更新一个组织的信息，用户必须是管理员
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                     true  "token"
// @Param        id       path      int                        true  "组织ID"
// @Param        data     body      model.UpdateOrganizationQ  true  "组织的名称，组织简介"
// @Success      200      {object}  model.CommonA              "是否成功，返回信息"
// @Router       /api/v1/organizations/{id} [put]
func UpdateOrganization(c *gin.Context) {
	var data model.CreateOrganizationQ
	if err := c.ShouldBindJSON(&data); err != nil {
		global.LOG.Panic("UpdateOrganization: bind data error")
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	organization, _ := service.GetOrganizationByID(id)
	if _, notFound := service.GetOrganizationByName(data.Name); !notFound && data.Name != organization.Name {
		global.LOG.Warn("UpdateOrganization: find same organization name")
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "已存在该名称的组织"})
		return
	}
	if organization, notFound := service.GetOrganizationByID(id); notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "未找到组织"})
	} else {
		service.UpdateOrganization(&organization, data.Name, data.Profile)
		c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新组织信息成功"})
	}
}

// DeleteOrganization
// @Summary      删除一个组织
// @Description  组织创建者删除该组织
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int            true  "组织ID"
// @Success      200      {object}  model.CommonA                   "是否成功，返回信息"
// @Router       /api/v1/organizations/{id} [delete]
func DeleteOrganization(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := service.DeleteOrganizationByID(id); err != nil {
		global.LOG.Panic("DeleteOrganization: delete organization error")
	}
	// TODO 删除已加入某组织的关系
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除组织成功"})
}

// CreateInvitation
// @Summary      邀请成员进入组织
// @Description  组织管理员向组织外人员发出邀请，邀请其成为管理员或普通用户
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int                      true  "组织ID"
// @Param        data     body      model.CreateInvitationQ  true  "用户email"
// @Success      200      {object}  model.CommonA            "是否成功，返回信息"
// @Router       /api/v1/organizations/{id}/invitations [post]
func CreateInvitation(c *gin.Context) {
	// 获取数据
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	data := utils.BindJsonData(c, &model.CreateInvitationQ{}).(*model.CreateInvitationQ)
	// 邀请成员存在性判定
	user, notFound := service.GetUserByEmail(data.Email)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "未找到被邀请用户"})
		return
	}
	// 组织存在性判定
	org, notFound := service.GetOrganizationByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "该组织不存在"})
		return
	}
	// 创建邀请
	err := service.CreateInvitation(&model.Invitation{
		UserID:   user.ID,
		UserName: user.Name,
		IsAdmin:  data.IsAdmin,
		OrgID:    id,
		OrgName:  org.Name})
	if err != nil {
		global.LOG.Panic("CreateInvitation: create invitation error")
	}
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "创建成功"})
}

// UpdateOrganizationMember
// @Summary      同意加入组织
// @Description  组织外的用户加入一个组织，该用户必须拥有组织管理员的邀请
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int            true  "组织ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/organizations/{id}/users [post]
func UpdateOrganizationMember(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	user := utils.SolveUser(c)
	if found, err := service.IsUserInThisOrganization(user.ID, id); found {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "用户已存在该组织中"})
	} else if err != nil {
		global.LOG.Panic("UpdateOrganizationMember: update invitation error: user or org not exist")
	}
	rel, _ := service.GetInValidInvitationByUserOrg(user.ID, id)
	rel.IsValid = true
	service.UpdateInvitation(rel)
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新成功"})
}

// GetOrganizationMember
// @Summary      获取组织成员
// @Description  获取一个组织中，所有组织成员的基础信息
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                        true  "token"
// @Param        id       path      int                             true  "组织ID"
// @Success      200      {object}  model.GetOrganizationMemberA  "是否成功，返回信息，组织成员信息列表"
// @Router       /api/v1/organizations/{id}/users [get]
func GetOrganizationMember(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if _, notFound := service.GetOrganizationByID(id); notFound {
		c.JSON(http.StatusOK, model.GetOrganizationMemberA{Success: false, Message: "找不到该组织的信息"})
	} else {
		c.JSON(http.StatusOK, model.GetOrganizationMemberA{Members: service.GetOrganizationMember(id), Success: true, Message: "获取成功"})
	}
}

// UpdateOrganizationAdmin
// @Summary      添加组织管理员
// @Description  组织创建者在组织中添加一个管理员，组织成员无法拒绝
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                          true  "token"
// @Param        id       path      int                           true  "组织ID"
// @Param        data     body      model.UpdateOrganizationAdminQ  true  "用户ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/organizations/{id}/admins [post]
func UpdateOrganizationAdmin(c *gin.Context) {
	user := utils.SolveUser(c)
	oid, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	data := utils.BindJsonData(c, &model.UpdateOrganizationAdminQ{}).(*model.UpdateOrganizationAdminQ)
	uid, _ := strconv.ParseUint(data.ID, 10, 64)
	organization, _ := service.GetOrganizationByID(oid)
	if organization.CreatorID != user.ID {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "组织成员无权操作"})
		return
	}
	if rel, notFound := service.GetInvitationByUserOrg(uid, oid); notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "成员未加入组织"})
	} else {
		rel.IsAdmin = true
		_ = service.UpdateInvitation(rel)
		c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新管理成功"})
	}
}

// DeleteOrganizationAdmin
// @Summary      删除组织管理员
// @Description  组织创建者在组织中删除一个管理员，管理员无法拒绝
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int            true  "组织ID"
// @Param        adminID  path      int            true  "管理员的用户ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/organizations/{id}/admins/{adminID} [delete]
func DeleteOrganizationAdmin(c *gin.Context) {
	user := utils.SolveUser(c)
	oid, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	uid, _ := strconv.ParseUint(c.Param("adminID"), 10, 64)
	organization, _ := service.GetOrganizationByID(oid)
	if organization.CreatorID != user.ID {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "组织成员无权操作"})
		return
	}
	if rel, notFound := service.GetInvitationByUserOrg(uid, oid); notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "成员未加入组织"})
	} else {
		rel.IsAdmin = false
		_ = service.UpdateInvitation(rel)
		c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除管理员成功"})
	}
}

// DeleteOrganizationMember
// @Summary      删除组织成员
// @Description  组织管理员在组织中删除一个成员，成员无法拒绝
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                   true  "token"
// @Param        id       path      int            true  "组织ID"
// @Param        userID   path      int            true  "用户ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/organizations/{id}/users/{userID} [delete]
func DeleteOrganizationMember(c *gin.Context) {
	// 获取请求数据
	user := utils.SolveUser(c)
	oid, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	uid, _ := strconv.ParseUint(c.Param("userID"), 10, 64)
	// 成员是否属于该组织
	rel, notFound := service.GetInvitationByUserOrg(uid, oid)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "成员未加入组织"})
		return
	}
	// 管理员权限判定
	for _, admin := range service.GetOrganizationAdmin(oid) {
		if admin == user.ID {
			global.DB.Delete(&rel)
			c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除成员成功"})
			return
		}
	}
	// 用户没有管理员权限
	c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "用户没有管理员权限"})
	return
}
