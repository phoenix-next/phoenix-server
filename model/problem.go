package model

import (
	"mime/multipart"
	"time"
)

type ResultT struct {
	ID          uint64    `json:"id"`
	Result      int       `json:"result"` // 0 AC , 1 WA , 2 TLE, 3 RE
	Language    string    `json:"language"`
	Path        string    `json:"path"`
	CreatedTime time.Time `json:"createdTime"`
}

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
	Name        string                `form:"name"`
	Difficulty  int                   `form:"difficulty"`
	Input       *multipart.FileHeader `form:"input" swaggerignore:"true"`
	Output      *multipart.FileHeader `form:"output" swaggerignore:"true"`
	Description *multipart.FileHeader `form:"description" swaggerignore:"true"`
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
	ProblemList []ProblemT `json:"problemList"` // 当前用户该题的评测结果，0 表示未做，1 表示通过，-1 表示评测过但是未通过
}

type UploadProblemRecordQ struct {
	Result   int                   `form:"result"` // 0 AC , 1 WA , 2 TLE, 3 RE
	Language string                `form:"language"`
	Code     *multipart.FileHeader `form:"code" swaggerignore:"true"`
}

type GetProblemRecordA struct {
	Success    bool      `json:"success"`
	Message    string    `json:"message"`
	ResultList []ResultT `json:"resultList"` // 0 AC , 1 WA , 2 TLE, 3 RE
}
