package model

type OrganizationT struct {
	OrgID   uint64 `json:"orgID"`
	OrgName string `json:"orgName"`
	IsAdmin bool   `json:"isAdmin"` // 用户在该组织中是否为管理员
}

type CommonA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type CreateUserQ struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

type CreateCaptchaQ struct {
	Email string `json:"email"`
}

type CreateTokenQ struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTokenA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
	ID      uint64 `json:"id"`
}

type GetUserA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

type GetUserOrganizationA struct {
	Success      bool            `json:"success"`
	Message      string          `json:"message"`
	Organization []OrganizationT `json:"organization"`
}

type GetUserInvitationA struct {
	Success      bool            `json:"success"`
	Message      string          `json:"message"`
	Organization []OrganizationT `json:"organization"`
}
