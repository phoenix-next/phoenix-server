package service

import (
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model/api"
	"github.com/phoenix-next/phoenix-server/model/database"
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
		return problem, err
	}
	return p, nil
}
