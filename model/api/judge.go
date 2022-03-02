package api

import "mime/multipart"

type CreateProblemQ struct {
	Name         string                `form:"name"`
	Difficulty   int                   `form:"difficulty"`
	Readable     int                   `form:"readable"`
	Writable     int                   `form:"writable"`
	Organization uint64                `form:"organization"`
	Input        *multipart.FileHeader `form:"input"`
	Output       *multipart.FileHeader `form:"output"`
	Description  *multipart.FileHeader `form:"description"`
}

type GetProblemQ struct {
	ID uint64 `json:"id"`
}

type GetProblemA struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Difficulty   int    `json:"difficulty"`
	Readable     int    `json:"readable"`
	Writable     int    `json:"writable"`
	Organization uint64 `json:"organization"`
	Input        string `json:"input"`
	Output       string `json:"output"`
	Description  string `json:"description"`
}

type UpdateProblemQ struct {
	ID           uint64                `form:"id"`
	Name         string                `form:"name"`
	Difficulty   int                   `form:"difficulty"`
	Readable     int                   `form:"readable"`
	Writable     int                   `form:"writable"`
	Organization uint64                `form:"organization"`
	Input        *multipart.FileHeader `form:"input"`
	Output       *multipart.FileHeader `form:"output"`
	Description  *multipart.FileHeader `form:"description"`
}

type DeleteProblemQ struct {
	ID uint64 `json:"id"`
}

type GetProblemVersionQ struct {
	ID uint64 `json:"id"`
}

type GetProblemVersionA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version int    `json:"version"`
}

type GetProblemListQ struct {
	Page   int `json:"page"`   //一页十个问题
	Sorter int `json:"sorter"` //按什么排序
}

type GetProblemListA struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ProblemList []struct {
		ID         uint64 `json:"id"`
		Name       string `json:"name"`
		Difficulty int    `json:"difficulty"`
	}
}
