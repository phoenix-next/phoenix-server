package api

import (
	"github.com/phoenix-next/phoenix-server/model/database"
	"mime/multipart"
)

type CreateProblemQ struct {
	Name         string                `form:"name"`
	Difficulty   int                   `form:"difficulty"`
	Readable     int                   `form:"readable"`
	Writable     int                   `form:"writable"`
	Creator      uint64                `form:"creator"`
	Organization uint64                `form:"organization"`
	Input        *multipart.FileHeader `form:"input" swaggerignore:"true"`
	Output       *multipart.FileHeader `form:"output" swaggerignore:"true"`
	Description  *multipart.FileHeader `form:"description" swaggerignore:"true"`
}

type GetProblemA struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Difficulty   int    `json:"difficulty"`
	Readable     int    `json:"readable"`
	Writable     int    `json:"writable"`
	Organization uint64 `json:"organization"`
	Creator      uint64 `json:"creator"`
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
	Input        *multipart.FileHeader `form:"input" swaggerignore:"true"`
	Output       *multipart.FileHeader `form:"output" swaggerignore:"true"`
	Description  *multipart.FileHeader `form:"description" swaggerignore:"true"`
}

type DeleteProblemQ struct {
	ID uint64 `json:"id"`
}

type GetProblemVersionA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version int    `json:"version"`
}

type GetProblemListA struct {
	Success     bool               `json:"success"`
	Message     string             `json:"message"`
	ProblemList []database.Problem `json:"problemList"`
}
