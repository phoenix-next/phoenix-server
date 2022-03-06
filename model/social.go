package model

type CommonA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RegisterQ struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

type GetCaptchaQ struct {
	Email string `json:"email"`
}

type LoginQ struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
	ID      uint64 `json:"id"`
}

type GetProfileA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

type GetUserOrganizationA struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	Organizations []struct {
		ID      uint64 `json:"id"`
		Name    string `json:"name"`
		IsAdmin bool   `json:"isAdmin"` // 用户在该组织中是否为管理员
	}
}

type CreateOrganizationQ struct {
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

type GetUserInvitationsA struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Organization []struct {
		ID      uint64 `json:"id"`
		Name    string `json:"name"`
		Profile string `json:"profile"`
	}
}

type GetAdminInfoA struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Organization []struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}
}
