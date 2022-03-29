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

// GetProblemFileFolder 获取保存某题目的文件夹
func GetProblemFileFolder(problemID uint64, version int) string {
	return strconv.Itoa(int(problemID)) + "_" + strconv.Itoa(version)
}

// GetProblemFileUrl 获取访问问题资源的Url
func GetProblemFileUrl(problem *model.Problem, kind string) string {
	return "resource/problem/" + GetProblemFileFolder(problem.ID, problem.Version) + "/" + kind
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
			res = strings.Compare(problems[i].Name, problems[j].Name) < 0
		case 3, -3:
			res = problems[i].Difficulty > problems[j].Difficulty
		}
		if sorter < 0 {
			res = !res
		}
		return
	})
	return problems[(page-1)*size : int(math.Min(float64(page*size), float64(len(problems))))]
}

// GetUserFinalJudge 返回单个题目的评测结果，0 表示未做，1 表示通过，-1 表示评测过但是未通过
func GetUserFinalJudge(uid uint64, pid uint64) int {
	results := QueryUserProblemResult(uid, pid)
	if len(results) == 0 {
		return 0
	}
	for _, r := range results {
		if r.Result == 0 {
			return 1
		}
	}
	return -1
}

// GetCodeFileName 根据记录获取保存的Code文件名称
func GetCodeFileName(result model.Result) string {
	// 关于时间格式化，可参考https://www.jianshu.com/p/c7f7fbb16932
	return strconv.FormatUint(result.ProblemID, 10) + "_" +
		strconv.FormatUint(result.UserID, 10) + "_" +
		result.CreatedTime.Format("20060102150405")
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
	problem.Name, problem.Version, problem.Difficulty = q.Name, problem.Version+1, q.Difficulty
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
	global.DB.Order("created_time desc").Find(&problems)
	return problems
}

// QueryUserProblemResult 查询用户对某问题的所有评测结果
func QueryUserProblemResult(uid uint64, pid uint64) (results []model.Result) {
	global.DB.Where("user_id = ? AND problem_id = ?", uid, pid).Find(&results)
	return results
}
