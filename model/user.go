package model

import "mime/multipart"

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
	Email   string `json:"email"`
	Profile string `json:"profile"`
	Avatar  string `json:"avatar"`
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

type UpdateUserQ struct {
	Name        string                `form:"name"`
	Password    string                `form:"password"`
	OldPassword string                `form:"oldPassword"`
	Profile     string                `form:"profile"`
	Avatar      *multipart.FileHeader `form:"avatar" swaggerignore:"true"`
}

type ForgetPasswordQ struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}
