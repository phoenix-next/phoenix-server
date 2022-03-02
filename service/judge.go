package service

import (
	"errors"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model/api"
	"github.com/phoenix-next/phoenix-server/model/database"
	"gorm.io/gorm"
	"strconv"
)

// MakeProblemFileName 保存题目文件名称
func MakeProblemFileName(problemId uint64, version int, suffix string) string {
	return strconv.Itoa(int(problemId)) + "_" + strconv.Itoa(version) + "_" + suffix
}

// CreateProblem 生成验证码
func CreateProblem(q *api.CreateProblemQ) (p database.Problem, err error) {
	//problem := database.Problem{Name: q.Name, Version: 1, Difficulty: 1, Readable: 1, Writable: 1, Organization: 1, Creator: 1}
	problem := database.Problem{Name: q.Name, Version: 1, Difficulty: q.Difficulty, Readable: q.Readable, Writable: q.Writable, Organization: q.Organization, Creator: q.Creator}
	if err = global.DB.Create(&problem).Error; err != nil {
		return p, err
	}
	return problem, nil
}

// GetProblemByID 根据问题 ID 查询某个问题
func GetProblemByID(ID uint64) (problem database.Problem, notFound bool) {
	err := global.DB.First(&problem, ID).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return problem, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Panic("GetUserByID: search error")
		return problem, true
	} else {
		return problem, false
	}
}

// 获取访问问题资源的文件名
func GetProblemFileName(problem *database.Problem, kind string) string {
	return strconv.Itoa(int(problem.ID)) + "_" + strconv.Itoa(problem.Version) + "_" + kind
}

// 获取访问问题资源的Url
func GetProblemFileUrl(problem *database.Problem, kind string) string {
	return "/resource/problem/" + GetProblemFileName(problem, kind)
}
