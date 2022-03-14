package model

import (
	"mime/multipart"
)

type CreateProblemQ struct {
	OrgID       uint64                `form:"organization"`
	Name        string                `form:"name"`
	Difficulty  int                   `form:"difficulty"`
	Readable    int                   `form:"readable"`
	Writable    int                   `form:"writable"`
	Input       *multipart.FileHeader `form:"input" swaggerignore:"true"`
	Output      *multipart.FileHeader `form:"output" swaggerignore:"true"`
	Description *multipart.FileHeader `form:"description" swaggerignore:"true"`
}

type GetProblemA struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	Name        string `json:"name"`
	Difficulty  int    `json:"difficulty"`
	Input       string `json:"input"`
	Output      string `json:"output"`
	Description string `json:"description"`
	Result      int    `json:"result"` // 当前用户该题的评测结果，0 表示未做，1 表示通过，-1 表示评测过但是未通过
}

type UpdateProblemQ struct {
	ID           uint64                `form:"id"`
	Name         string                `form:"name"`
	Difficulty   int                   `form:"difficulty"`
	Readable     int                   `form:"readable"`
	Writable     int                   `form:"writable"`
	Organization uint64                `form:"organization"`
	Input        *multipart.FileHeader `form:"input" swaggerignore:"true"`
	Output       *multipart.FileHeader `form:"output" swaggerignore:"true"`
	Description  *multipart.FileHeader `form:"description" swaggerignore:"true"`
}

type GetProblemVersionA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version int    `json:"version"`
}

type GetProblemListA struct {
	Success     bool       `json:"success"`
	Message     string     `json:"message"`
	Total       int        `json:"total"`
	ProblemList []ProblemT `json:"problemList"`
}

type UploadProblemRecordQ struct {
	Result int                   `form:"result"`
	Code   *multipart.FileHeader `form:"code" swaggerignore:"true"`
}

type GetProblemRecordA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
