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
// @Param        x-token  header    string                  true  "token"
// @Param        data     body      model.CreateOrganizationQ  true  "组织名称，组织的简介"
// @Success      200      {object}  model.CommonA              "是否成功，返回信息"
// @Router       /api/v1/organizations [post]
func CreateOrganization(c *gin.Context) {
	// 获取请求数据
	data := utils.BindJsonData(c, &model.CreateOrganizationQ{}).(*model.CreateOrganizationQ)
	user := utils.SolveUser(c)
	// 有重名组织的情况
	if _, notFound := service.GetOrganizationByName(data.Name); !notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "已存在该名称的组织"})
		return
	}
	// 创建组织，保存组织元数据
	org := model.Organization{
		Name:        data.Name,
		Profile:     data.Profile,
		CreatorName: user.Name,
		CreatorID:   user.ID}
	global.DB.Create(&org)
	// 维护用户 - 组织关系
	global.DB.Create(&model.Invitation{
		UserID:    user.ID,
		UserName:  user.Name,
		UserEmail: user.Email,
		OrgID:     org.ID,
		OrgName:   org.Name,
		IsAdmin:   true,
		IsValid:   true})
	// 返回响应
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "创建组织成功"})
}

// GetOrganization
// @Summary      获取组织信息
// @Description  获取一个组织的详细信息
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                     true  "token"
// @Param        id       path      int                     true  "组织ID"
// @Success      200      {object}  model.GetOrganizationA  "是否成功，返回信息，组织名称，组织简介，当前用户是否在组织中，当前用户是否为管理员"
// @Router       /api/v1/organizations/{id} [get]
func GetOrganization(c *gin.Context) {
	// 获取请求数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.GetOrganizationA{Success: false, Message: "组织ID不合法"})
	}
	// 组织存在性判定
	org, notFound := service.GetOrganizationByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.GetOrganizationA{Success: false, Message: "找不到该组织"})
		return
	}
	// 获取当前用户与组织的关系
	user := utils.SolveUser(c)
	rel, notFound := service.GetInvitationByUserOrg(user.ID, org.ID)
	// 当前用户不是组织中的成员
	if notFound {
		c.JSON(http.StatusOK, model.GetOrganizationA{
			Success: true,
			Message: "获取组织信息成功",
			Name:    org.Name,
			Profile: org.Profile,
			IsValid: false,
			IsAdmin: false})
		return
	}
	// 当前用户是组织中的成员
	c.JSON(http.StatusOK, model.GetOrganizationA{
		Success: true,
		Message: "获取组织信息成功",
		Name:    org.Name,
		Profile: org.Profile,
		IsValid: true,
		IsAdmin: rel.IsAdmin})

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
	// 获取请求数据
	data := utils.BindJsonData(c, &model.CreateOrganizationQ{}).(*model.CreateOrganizationQ)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "组织ID非法"})
		return
	}
	// 组织的存在性判定
	org, notFound := service.GetOrganizationByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "该组织不存在"})
		return
	}
	// 组织重名的情况
	if _, notFound = service.GetOrganizationByName(data.Name); !notFound && data.Name != org.Name {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "已存在该名称的组织"})
		return
	}
	// 成功更新信息
	org.Name = data.Name
	org.Profile = data.Profile
	global.DB.Save(&org)
	// 维护成员 - 组织关系
	global.DB.Model(&model.Invitation{}).Where("org_id = ?", org.ID).Update("org_name", org.Name)
	// 返回结果
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新组织信息成功"})

}

// DeleteOrganization
// @Summary      解散一个组织
// @Description  组织创建者解散该组织
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int            true  "组织ID"
// @Success      200      {object}  model.CommonA                   "是否成功，返回信息"
// @Router       /api/v1/organizations/{id} [delete]
func DeleteOrganization(c *gin.Context) {
	// 获取请求数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "请求参数非法"})
		return
	}
	// 组织的存在性判定
	org, notFound := service.GetOrganizationByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "该组织不存在"})
		return
	}
	// 用户权限判定
	user := utils.SolveUser(c)
	if user.ID != org.CreatorID {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "用户没有权限进行该操作"})
		return
	}
	// 维护成员 - 组织关系
	global.DB.Where("org_id = ?", org.ID).Delete(&model.Invitation{})
	// 删除组织元数据
	global.DB.Delete(&org)
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "解散组织成功"})
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
	// 获取请求数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	data := utils.BindJsonData(c, &model.CreateInvitationQ{}).(*model.CreateInvitationQ)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "请求参数非法"})
		return
	}
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
	// 无需重复创建邀请
	var invitation []model.Invitation
	global.DB.Model(&model.Invitation{}).Where("user_id = ? AND org_id = ?", user.ID, org.ID).Find(&invitation)
	if len(invitation) > 0 {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "已经邀请过该用户"})
		return
	}
	// 创建邀请
	err = service.CreateInvitation(&model.Invitation{
		UserID:    user.ID,
		UserName:  user.Name,
		UserEmail: user.Email,
		IsAdmin:   data.IsAdmin,
		OrgID:     id,
		OrgName:   org.Name})
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
	// 获取请求数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "组织ID非法"})
		return
	}
	user := utils.SolveUser(c)
	// 用户与组织的关系判定
	found, err := service.IsUserInThisOrganization(user.ID, id)
	if found {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "用户已存在该组织中"})
		return
	}
	// 组织不存在的情况
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "该组织不存在"})
		return
	}
	// 更新数据库
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
	// 获取请求数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "组织ID非法"})
		return
	}
	// 找不到组织的情况
	if _, notFound := service.GetOrganizationByID(id); notFound {
		c.JSON(http.StatusOK, model.GetOrganizationMemberA{Success: false, Message: "找不到该组织的信息"})
		return
	}
	// 成功返回
	c.JSON(http.StatusOK, model.GetOrganizationMemberA{
		Members: service.GetOrganizationMember(id),
		Success: true,
		Message: "获取组织成员成功"})
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
	// 获取请求参数
	oid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "组织ID非法"})
		return
	}
	user := utils.SolveUser(c)
	uid := utils.BindJsonData(c, &model.UpdateOrganizationAdminQ{}).(*model.UpdateOrganizationAdminQ).ID
	// 组织不存在的情况
	org, notFound := service.GetOrganizationByID(oid)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "该组织不存在"})
		return
	}
	// 请求用户权限判定
	if org.CreatorID != user.ID {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "非组织创建者无权进行该操作"})
		return
	}
	// 成员未加入组织的情况
	rel, notFound := service.GetInvitationByUserOrg(uid, oid)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "成员未加入组织"})
		return
	}
	// 成功返回
	rel.IsAdmin = true
	_ = service.UpdateInvitation(rel)
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "添加管理员成功"})

}

// DeleteOrganizationAdmin
// @Summary      取消组织管理员
// @Description  组织创建者在组织中取消某管理员的管理员权限，管理员无法拒绝
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int            true  "组织ID"
// @Param        adminID  path      int            true  "管理员的用户ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/organizations/{id}/admins/{adminID} [delete]
func DeleteOrganizationAdmin(c *gin.Context) {
	// 获取请求中的数据
	user := utils.SolveUser(c)
	oid, err1 := strconv.ParseUint(c.Param("id"), 10, 64)
	uid, err2 := strconv.ParseUint(c.Param("adminID"), 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "请求参数非法"})
		return
	}
	// 组织的存在性判定
	org, notFound := service.GetOrganizationByID(oid)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "该组织不存在"})
		return
	}
	// 当前用户权限判定
	if org.CreatorID != user.ID {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "组织成员无权操作"})
		return
	}
	// 创建者不能取消自己的管理员权限
	if org.CreatorID == uid {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "您无法取消自己的管理员权限"})
		return
	}
	// 操作对象不在组织中的情况
	rel, notFound := service.GetInvitationByUserOrg(uid, oid)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "成员未加入组织"})
		return
	}
	// 请求成功
	rel.IsAdmin = false
	_ = service.UpdateInvitation(rel)
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "取消管理员成功"})

}

// DeleteOrganizationMember
// @Summary      踢出组织成员
// @Description  组织管理员在组织中删除一个非管理员成员，该成员无法拒绝
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
	oid, err1 := strconv.ParseUint(c.Param("id"), 10, 64)
	uid, err2 := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "请求参数非法"})
		return
	}
	// 成员不属于该组织的情况
	rel, notFound := service.GetInvitationByUserOrg(uid, oid)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "成员未加入组织"})
		return
	}
	// 管理员无法踢出一个管理员
	invitation := service.GetOrganizationAdmin(oid)
	for _, admin := range invitation {
		if admin.UserID == uid {
			c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "无法删除具有管理员权限的成员"})
			return
		}
	}
	// 管理员将成员踢出组织
	for _, admin := range invitation {
		if admin.UserID == user.ID {
			global.DB.Delete(&rel)
			c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除成员成功"})
			return
		}
	}
	// 用户没有管理员权限
	c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "用户没有管理员权限"})
	return
}
