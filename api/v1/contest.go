package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"github.com/phoenix-next/phoenix-server/service"
	"github.com/phoenix-next/phoenix-server/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// CreateContest
// @Summary      创建比赛
// @Description  组织管理员创建一个比赛，比赛包含指定题号的题目
// @Tags         比赛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        data     body      model.CreateContestQ  true  "组织ID，比赛名称，比赛简介，可读权限，开始时间，结束时间，题目列表"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/contests [post]
func CreateContest(c *gin.Context) {
	// 获取请求数据
	data := utils.BindJsonData(c, &model.CreateContestQ{}).(*model.CreateContestQ)
	user := utils.SolveUser(c)
	// 用户权限判定
	for _, admin := range service.GetOrganizationAdmin(data.OrgID) {
		if admin.UserID == user.ID {
			// 创建比赛
			contest := &model.Contest{
				OrgID:     data.OrgID,
				Profile:   data.Profile,
				Name:      data.Name,
				Readable:  data.Readable,
				StartTime: data.StartTime,
				EndTime:   data.EndTime}
			global.DB.Create(&contest)
			// 维护比赛 - 题目关系
			for _, problemID := range data.ProblemIDs {
				// 获得对应ID的题目
				problem, notFound := service.GetProblemByID(problemID)
				if notFound {
					c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "题目列表中有题目不存在"})
					return
				}
				// 维护关系
				global.DB.Create(&model.ContestProblem{
					ContestID:   contest.ID,
					ProblemID:   problem.ID,
					ProblemName: problem.Name,
					Difficulty:  problem.Difficulty})
			}
			// 返回结果
			c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "创建比赛成功"})
			return
		}
	}
	// 用户没有权限
	c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "仅有组织管理员才能创建比赛"})
}

// GetContest
// @Summary      获取比赛信息
// @Description  获取一个比赛的详细信息，包括该比赛的名称以及包含题目等信息
// @Tags         比赛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       path      int            true  "比赛ID"
// @Success      200      {object}  model.GetContestA  "是否成功，返回信息，比赛名称，比赛简介，题目列表"
// @Router       /api/v1/contests/{id} [get]
func GetContest(c *gin.Context) {
	// 获取请求数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.GetContestA{Success: false, Message: "比赛ID不合法"})
		return
	}
	// 比赛的存在性判定
	var contest model.Contest
	if err := global.DB.First(&contest, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, model.GetContestA{Success: false, Message: "比赛不存在"})
		return
	}
	// 获取所有题目
	var problems []model.ProblemT
	global.DB.Model(&model.ContestProblem{}).Where("contest_id = ?", contest.ID).Find(&problems)
	// 返回结果
	c.JSON(http.StatusOK, model.GetContestA{
		Success: true,
		Message: "获取比赛信息成功",
		Name:    contest.Name,
		Profile: contest.Profile,
		Problem: problems})
}

// UpdateContest
// @Summary      更新比赛
// @Description  管理员更新一个比赛的信息，如更新比赛包含的题目、比赛名称等信息
// @Tags         比赛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                         true  "token"
// @Param        data     body      model.UpdateContestQ  true  "比赛名称，比赛简介，比赛包含的题目ID"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/contests/{id} [put]
func UpdateContest(c *gin.Context) {
	// 获取请求数据
	data := utils.BindJsonData(c, &model.UpdateContestQ{}).(*model.UpdateContestQ)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "比赛ID不合法"})
		return
	}
	// 比赛的存在性判定
	var contest model.Contest
	if err := global.DB.First(&contest, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "比赛不存在"})
		return
	}
	// 修改比赛的元数据
	contest.Name = data.Name
	contest.Profile = data.Profile
	global.DB.Save(&contest)
	// 更新题目 - 比赛关系
	global.DB.Where("contest_id = ?", contest.ID).Delete(&model.ContestProblem{})
	for _, problemID := range data.ProblemIDs {
		// 获得对应ID的题目
		problem, notFound := service.GetProblemByID(problemID)
		if notFound {
			c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "题目列表中有题目不存在"})
			return
		}
		// 维护关系
		global.DB.Create(&model.ContestProblem{
			ContestID:   contest.ID,
			ProblemID:   problem.ID,
			ProblemName: problem.Name,
			Difficulty:  problem.Difficulty})
	}
	// 返回响应
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "修改比赛信息成功"})
}

// DeleteContest
// @Summary      删除比赛
// @Description  管理员永久删除一个比赛
// @Tags         比赛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        id       path      int                true  "比赛ID"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/contests/{id} [delete]
func DeleteContest(c *gin.Context) {
	// 获取请求数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "比赛ID不合法"})
		return
	}
	// 比赛的存在性判定
	var contest model.Contest
	if err := global.DB.First(&contest, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "比赛不存在"})
		return
	}
	// 用户权限判定
	user := utils.SolveUser(c)
	for _, admin := range service.GetOrganizationAdmin(contest.OrgID) {
		if user.ID == admin.UserID {
			// 维护问题 - 比赛关系
			global.DB.Where("contest_id = ?", contest.ID).Delete(&model.ContestProblem{})
			// 删除比赛元数据
			global.DB.Delete(&contest)
			// 返回响应
			c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除比赛成功"})
			return
		}
	}
	// 用户没有权限操作
	c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "用户没有管理员权限"})
}

// GetContestList
// @Summary      获取比赛列表
// @Description  获取用户所能查看的比赛列表
// @Tags         比赛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                 true  "token"
// @Param        page     query     int                    true  "用户位于哪一页，页数从1开始"
// @Param        keyWord  query     string                 true  "当前的(比赛名称)搜索关键字，为空字符串表示没有关键字，模糊匹配"
// @Param        sorter   query     int                    true  "用户想按什么排序，1为按ID升序，-1为按ID降序，2为按名称升序，-2为按名称降序"
// @Success      200      {object}  model.GetContestListA  "是否成功，返回信息，比赛列表"
// @Router       /api/v1/contests [get]
func GetContestList(c *gin.Context) {
	// 获取请求数据
	page, err1 := strconv.Atoi(c.Query("page"))
	sorter, err2 := strconv.Atoi(c.Query("sorter"))
	keyWord := c.Query("keyWord")
	user := utils.SolveUser(c)
	// 请求数据不合法的情况
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, model.GetContestListA{Success: false, Message: "请求参数不合法"})
		return
	}
	// 获取可读的比赛,并进行排序
	contests := service.GetReadableContest(user.ID, sorter)
	// 对比赛标题进行模糊查找
	var filteredContests []model.ContestT
	for _, contest := range contests {
		if fuzzy.Match(keyWord, contest.Name) {
			filteredContests = append(filteredContests, contest)
		}
	}
	// 得到比赛的总页数
	totalPage := len(filteredContests) / 5
	if len(filteredContests)%5 != 0 {
		totalPage += 1
	}
	// 查不到比赛的情况
	if totalPage == 0 {
		c.JSON(http.StatusOK, model.GetContestListA{
			Success:     true,
			Message:     "获取比赛列表成功",
			Total:       0,
			ContestList: []model.ContestT{}})
		return
	}
	// 页数不合法的情况
	if page <= 0 || page > totalPage {
		c.JSON(http.StatusOK, model.GetContestListA{Success: false, Message: "页数非法"})
		return
	}
	// 获取端点位置，并对帖子切片
	start, end := (page-1)*5, page*5
	if length := len(filteredContests); end > length {
		end = length
	}
	slicedContests := filteredContests[start:end]
	// 返回
	c.JSON(http.StatusOK, model.GetContestListA{
		Success:     true,
		Message:     "获取比赛列表成功",
		Total:       len(filteredContests),
		ContestList: slicedContests})
}

// GetOrganizationProblem
// @Summary      获取组织中的题目
// @Description  获取组织中的管理员可见的题目，即属于组织管理员可读的题目
// @Tags         比赛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        id       path      int                            true  "组织ID"
// @Success      200      {object}  model.GetOrganizationProblemA  "是否成功，返回信息，题目列表"
// @Router       /api/v1/organizations/{id}/problems [get]
func GetOrganizationProblem(c *gin.Context) {
	// 获取请求数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.GetOrganizationProblemA{Success: false, Message: "组织ID不合法"})
		return
	}
	// 组织的存在性判定
	_, notFound := service.GetOrganizationByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.GetOrganizationProblemA{Success: false, Message: "该组织不存在"})
		return
	}
	// 用户权限判定
	user := utils.SolveUser(c)
	for _, admin := range service.GetOrganizationAdmin(id) {
		if user.ID == admin.UserID {
			// 获取组织中的题目
			var problems []model.Problem
			global.DB.Where("org_id = ? AND (readable = ? OR readable = ?)", id, 2, 1).Find(&problems)
			// 返回响应
			c.JSON(http.StatusOK, model.GetOrganizationProblemA{Success: true, Message: "获取组织可见的题目成功", ProblemList: problems})
			return
		}
	}
	// 用户没有权限操作
	c.JSON(http.StatusOK, model.GetOrganizationProblemA{Success: false, Message: "用户没有管理员权限"})
}
