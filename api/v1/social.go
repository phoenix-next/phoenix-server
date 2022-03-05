package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateOrganization
// @Summary      创建一个组织
// @Description  创建一个组织，创建者比管理员权限高，也算是管理员
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string       true  "token"
// @Param        data     body      api.CreateOrganizationQ  true  "组织名称，组织的简介"
// @Success      200      {object}  api.CommonA              "是否成功，返回信息"
// @Router       /api/v1/organizations [post]
func CreateOrganization(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// UpdateOrganization
// @Summary      更新一个组织的信息
// @Description  更新一个组织的信息，用户必须是管理员
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                   true  "token"
// @Param        id       path      int                      true  "组织ID"
// @Param        data     body      api.UpdateOrganizationQ  true  "组织的名称，组织简介"
// @Success      200      {object}  api.CommonA              "是否成功，返回信息"
// @Router       /api/v1/organizations/{id} [put]
func UpdateOrganization(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// DeleteOrganization
// @Summary      删除一个组织
// @Description  组织创建者删除该组织
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                        true  "token"
// @Param        id       path      int          true  "组织ID"
// @Success      200      {object}  api.CommonA                   "是否成功，返回信息"
// @Router       /api/v1/organizations/{id} [delete]
func DeleteOrganization(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// CreateInvitation
// @Summary      邀请成员进入组织
// @Description  组织管理员向组织外人员发出邀请
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                 true  "token"
// @Param        id       path      int                    true  "组织ID"
// @Param        data     body      api.CreateInvitationQ  true  "用户email"
// @Success      200      {object}  api.CommonA            "是否成功，返回信息"
// @Router       /api/v1/organizations/{id}/invitations [post]
func CreateInvitation(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// UpdateOrganizationMember
// @Summary      同意加入组织
// @Description  组织外的用户加入一个组织，该用户必须拥有组织管理员的邀请
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string       true  "token"
// @Param        id       path      int          true  "组织ID"
// @Success      200      {object}  api.CommonA  "是否成功，返回信息"
// @Router       /api/v1/organizations/{id}/users [post]
func UpdateOrganizationMember(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// GetOrganizationMember
// @Summary      获取组织成员
// @Description  获取一个组织中，所有组织成员的基础信息
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                      true  "token"
// @Param        id       path      int                         true  "组织ID"
// @Success      200      {object}  api.GetOrganizationMemberA  "是否成功，返回信息，组织成员信息列表"
// @Router       /api/v1/organizations/{id}/users [get]
func GetOrganizationMember(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// UpdateOrganizationAdmin
// @Summary      添加组织管理员
// @Description  组织创建者在组织中添加一个管理员，组织成员无法拒绝
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string       true  "token"
// @Param        id       path      int                           true  "组织ID"
// @Param        data     body      api.UpdateOrganizationAdminQ  true  "用户ID"
// @Success      200      {object}  api.CommonA  "是否成功，返回信息"
// @Router       /api/v1/organizations/{id}/admins [post]
func UpdateOrganizationAdmin(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// DeleteOrganizationAdmin
// @Summary      删除组织管理员
// @Description  组织创建者在组织中删除一个管理员，管理员无法拒绝
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                   true  "token"
// @Param        id       path      int          true  "组织ID"
// @Param        adminID  path      int          true  "管理员的用户ID"
// @Success      200      {object}  api.CommonA  "是否成功，返回信息"
// @Router       /api/v1/organizations/{id}/admins/{adminID} [delete]
func DeleteOrganizationAdmin(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}
