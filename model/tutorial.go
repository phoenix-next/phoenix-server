package model

import "mime/multipart"

type TutorialT struct {
	ID          uint64 `json:"id"`
	CreatorName string `json:"creatorName"`
	Name        string `json:"name"`
	Profile     string `json:"profile"`
}

type CreateTutorialQ struct {
	OrgID    uint64                `form:"orgID"`
	Name     string                `form:"name"`
	Profile  string                `form:"profile"`
	Readable int                   `form:"readable"`
	Writable int                   `form:"writable"`
	File     *multipart.FileHeader `form:"file" swaggerignore:"true"`
}

type GetTutorialA struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	OrgID        uint64 `json:"orgID"`
	CreatorID    uint64 `json:"creatorID"`
	CreatorName  string `json:"creatorName"`
	Name         string `json:"name"`
	Profile      string `json:"profile"`
	Version      int    `json:"version"`
	TutorialPath string `json:"tutorialPath"` // 教程的下载路径
}

type UpdateTutorialQ struct {
	ID      uint64                `form:"id"`
	Name    string                `form:"name"`
	Profile string                `form:"profile"`
	File    *multipart.FileHeader `form:"file" swaggerignore:"true"`
}

type GetTutorialVersionA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version int    `json:"version"`
}

type GetTutorialListA struct {
	Success      bool        `json:"success"`
	Message      string      `json:"message"`
	Total        int         `json:"total"`
	TutorialList []TutorialT `json:"tutorialList"`
}
