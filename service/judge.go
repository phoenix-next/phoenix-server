package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model/api"
	"github.com/phoenix-next/phoenix-server/model/database"
	"gorm.io/gorm"
	"math"
	"sort"
	"strconv"
	"strings"
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
	return "resource/problem/" + GetProblemFileName(problem, kind)
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

// QueryAllProblems 查询所有问题
func QueryAllProblems() (problems []database.Problem) {
	problems = make([]database.Problem, 0)
	global.DB.Find(&problems)
	return problems
}

// GetAllAvailableReadableProblems 获取所有可访问问题 TODO 组织管理员权限
func GetAllAvailableReadableProblems(c *gin.Context) (problems []database.Problem) {
	allProblems := QueryAllProblems()
	problems = make([]database.Problem, 0)
	for _, problem := range allProblems {
		if problem.Readable == 3 {
			problems = append(problems, problem)
		}
	}
	if _, isUser := c.Get("email"); !isUser {
		// 未登录，直接返回公开题目
		return problems
	}
	user, _ := GetUserByEmail(c.GetString("email"))
	//TODO 组织管理员权限
	for _, problem := range allProblems {
		if problem.Creator == user.ID {
			problems = append(problems, problem)
		}
	}
	return problems
}

// GetProblemsByPage 对给定问题做出排序与页选择
func GetProblemsByPage(problems []database.Problem, page int, sorter int) (problemList []database.Problem) {
	size := 10
	sort.Slice(problems, func(i, j int) bool {
		if math.Abs(float64(sorter)) == 1 {
			return problems[i].ID > problems[j].ID && sorter > 0
		} else if math.Abs(float64(sorter)) == 3 {
			return problems[i].Difficulty > problems[j].Difficulty && sorter > 0
		} else if math.Abs(float64(sorter)) == 2 {
			return strings.Compare(problems[i].Name, problems[j].Name) < 0 && sorter > 0
		} else {
			return true
		}
	})
	return problems[(page-1)*size : int(math.Min(float64(page*size), float64(len(problems))))]
}
