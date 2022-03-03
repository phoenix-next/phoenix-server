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

// GetProblemFileName 获取访问问题资源的文件名
func GetProblemFileName(problem *database.Problem, kind string) string {
	return strconv.Itoa(int(problem.ID)) + "_" + strconv.Itoa(problem.Version) + "_" + kind
}

// GetProblemFileUrl 获取访问问题资源的Url
func GetProblemFileUrl(problem *database.Problem, kind string) string {
	return "/resource/problem/" + GetProblemFileName(problem, kind)
}

// DeleteProblemByID 根据问题id 删除问题
func DeleteProblemByID(ID uint64) (err error) {
	if err = global.DB.Where("id = ?", ID).Delete(database.Problem{}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateProblem 根据信息更新题目
func UpdateProblem(problem *database.Problem, q *api.UpdateProblemQ) (err error) {
	problem.Name, problem.Version, problem.Difficulty, problem.Organization, problem.Writable, problem.Readable =
		q.Name, problem.Version+1, q.Difficulty, q.Organization, q.Writable, q.Readable
	err = global.DB.Save(problem).Error
	return err
}

// SaveProblem 根据信息保存题目
func SaveProblem(problem *database.Problem) (err error) {
	err = global.DB.Save(problem).Error
	return err
}
