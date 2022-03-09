package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"github.com/phoenix-next/phoenix-server/utils"
	"gorm.io/gorm"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Helper

// MakeProblemFileName 保存题目文件名称
func MakeProblemFileName(problemId uint64, version int, suffix string) string {
	return strconv.Itoa(int(problemId)) + "_" + strconv.Itoa(version) + "_" + suffix
}

// GetProblemFileName 获取访问问题资源的文件名
func GetProblemFileName(problem *model.Problem, kind string) string {
	return strconv.Itoa(int(problem.ID)) + "_" + strconv.Itoa(problem.Version) + "_" + kind
}

// GetProblemFileUrl 获取访问问题资源的Url
func GetProblemFileUrl(problem *model.Problem, kind string) string {
	return "resource/problem/" + GetProblemFileName(problem, kind)
}

// GetReadableProblems 获取所有可访问问题
func GetReadableProblems(c *gin.Context) (problems []model.Problem) {
	user := utils.SolveUser(c)
	allProblems := QueryAllProblems()
	problems = make([]model.Problem, 0)

	// 创建组织到是否是组织管理员的映射
	orgAdminMap := make(map[uint64]bool)
	for _, invitation := range GetUserOrganization(user.ID) {
		orgAdminMap[invitation.OrgID] = invitation.IsAdmin
	}
	for _, problem := range allProblems {
		if problem.Readable == 3 || problem.Creator == user.ID {
			problems = append(problems, problem)
		} else if isAdmin, ok := orgAdminMap[problem.OrgID]; ok && isAdmin {
			// 在题目的组织中且是管理员，直接加入
			problems = append(problems, problem)
		} else if ok && !isAdmin && problem.Readable == 2 {
			// 在题目所属组织中，但不是管理员，仅当可读为2
			problems = append(problems, problem)
		}
	}

	return problems
}

// GetProblemsByPage 对给定问题做出排序与页选择
func GetProblemsByPage(problems []model.Problem, page int, sorter int) (problemList []model.Problem) {
	size := 10
	sort.Slice(problems, func(i, j int) (res bool) {
		switch sorter {
		case 1, -1:
			res = problems[i].ID > problems[j].ID
		case 2, -2:
			res = problems[i].Difficulty > problems[j].Difficulty
		case 3, -3:
			res = strings.Compare(problems[i].Name, problems[j].Name) < 0
		}
		if sorter < 0 {
			res = !res
		}
		return
	})
	return problems[(page-1)*size : int(math.Min(float64(page*size), float64(len(problems))))]
}

// 数据库操作

// GetProblemByID 根据问题 ID 查询某个问题
func GetProblemByID(ID uint64) (problem model.Problem, notFound bool) {
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

// DeleteProblemByID 根据问题id 删除问题
func DeleteProblemByID(ID uint64) (err error) {
	if err = global.DB.Where("id = ?", ID).Delete(model.Problem{}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateProblem 根据信息更新题目
func UpdateProblem(problem *model.Problem, q *model.UpdateProblemQ) (err error) {
	problem.Name, problem.Version, problem.Difficulty, problem.OrgID, problem.Writable, problem.Readable =
		q.Name, problem.Version+1, q.Difficulty, q.Organization, q.Writable, q.Readable
	err = global.DB.Save(problem).Error
	return err
}

// SaveProblem 根据信息保存题目
func SaveProblem(problem *model.Problem) (err error) {
	err = global.DB.Save(problem).Error
	return err
}

// QueryAllProblems 查询所有问题
func QueryAllProblems() (problems []model.Problem) {
	problems = make([]model.Problem, 0)
	global.DB.Find(&problems)
	return problems
}
