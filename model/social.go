package model

type CreateOrganizationQ struct {
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

type GetOrganizationA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

type UpdateOrganizationQ struct {
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

type CreateInvitationQ struct {
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
}

type Member struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"` // 用户在该组织中是否为管理员
}

type GetOrganizationMemberA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Members []Member
}

type UpdateOrganizationAdminQ struct {
	ID string `json:"id"`
}
