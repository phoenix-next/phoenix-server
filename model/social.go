package model

type CreateOrganizationQ struct {
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

type GetOrganizationA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Profile string `json:"profile"`
	IsValid bool   `json:"isValid"`
	IsAdmin bool   `json:"isAdmin"`
}

type UpdateOrganizationQ struct {
	Name    string `form:"name"`
	Profile string `form:"profile"`
}

type CreateInvitationQ struct {
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
}

type Member struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"` // 用户在该组织中是否为管理员
}

type GetOrganizationMemberA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Members []Member
}

type UpdateOrganizationAdminQ struct {
	ID uint64 `json:"id"`
}

type UpdateOrganizationMemberQ struct {
	Accept bool `json:"accept"`
}
