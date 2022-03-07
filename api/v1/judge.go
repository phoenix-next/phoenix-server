package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"github.com/phoenix-next/phoenix-server/service"
	"github.com/phoenix-next/phoenix-server/utils"
	"net/http"
	"path/filepath"
	"strconv"
)

// CreateProblem
// @Summary      创建题目
// @Description  创建一个题目，题目需要包含答案和题面
// @Tags         评测模块
// @Accept       multipart/form-data
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        data     body      model.CreateProblemQ  true  "题目名称，题目难度，可读权限，可写权限，组织ID，输入文件，输出文件，题目描述"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/problems [post]
func CreateProblem(c *gin.Context) {
	// 获取题目保存路径，获取用户
	path := global.VP.GetString("problem_path")
	user := utils.SolveUser(c)
	// 获取请求数据
	var data model.CreateProblemQ
	if c.ShouldBind(&data) != nil {
		global.LOG.Panic("CreateProblem: bind data error")
	}
	// 创建题目
	problem := model.Problem{
		Name:         data.Name,
		Version:      1,
		Difficulty:   data.Difficulty,
		Readable:     data.Readable,
		Writable:     data.Writable,
		Organization: data.Organization,
		Creator:      user.ID}
	if global.DB.Create(&problem).Error != nil {
		global.LOG.Warn("CreateProblem: create problem error")
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "创建题目失败"})
		return
	}
	// 保存题目相关的文件
	err1 := c.SaveUploadedFile(data.Description, filepath.Join(path, service.MakeProblemFileName(problem.ID, 1, "description")))
	err2 := c.SaveUploadedFile(data.Input, filepath.Join(path, service.MakeProblemFileName(problem.ID, 1, "input")))
	err3 := c.SaveUploadedFile(data.Output, filepath.Join(path, service.MakeProblemFileName(problem.ID, 1, "output")))
	if err1 != nil || err2 != nil || err3 != nil {
		//发生错误，回滚数据库
		service.DeleteProblemByID(problem.ID)
		global.LOG.Panic("CreateProblem: save problem error")
	}
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "创建题目成功"})
}

// GetProblem
// @Summary      下载题目
// @Description  下载一个题目的信息
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       path      int            true  "题目ID"
// @Success      200      {object}  model.GetProblemA  "题目ID，题目名称，题目难度，可读权限，可写权限，组织ID，输入文件，输出文件，题目描述"
// @Router       /api/v1/problems/{id} [get]
func GetProblem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if problem, notFound := service.GetProblemByID(id); notFound {
		c.JSON(http.StatusOK, model.GetProblemA{Success: false, Message: "找不到该题目的信息"})
	} else {
		c.JSON(http.StatusOK, model.GetProblemA{
			Success:      true,
			Message:      "获取题目成功",
			ID:           problem.ID,
			Name:         problem.Name,
			Difficulty:   problem.Difficulty,
			Readable:     problem.Readable,
			Writable:     problem.Writable,
			Organization: problem.Organization,
			Creator:      problem.Creator,
			Input:        service.GetProblemFileUrl(&problem, "input"),
			Output:       service.GetProblemFileUrl(&problem, "output"),
			Description:  service.GetProblemFileUrl(&problem, "description"),
		})
	}
}

// UpdateProblem
// @Summary      更新题目
// @Description  更新一个题目的信息，并自动更新题目版本
// @Tags         评测模块
// @Accept       multipart/form-data
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        data     body      model.UpdateProblemQ  true  "题目ID，题目名称，题目难度，可读权限，可写权限，组织ID，输入文件，输出文件，题目描述"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/problems/{id} [put]
func UpdateProblem(c *gin.Context) {
	var data model.UpdateProblemQ
	path := filepath.Join(global.VP.GetString("root_path"), "resource", "problems")
	err := c.ShouldBind(&data)
	if err != nil {
		global.LOG.Panic("UpdateProblem: bind data error")
	}

	if problem, notFound := service.GetProblemByID(data.ID); notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "找不到该题目的信息"})
	} else {
		problemOrigin := problem
		err = service.UpdateProblem(&problem, &data)
		if err != nil {
			global.LOG.Panic("UpdateProblem: save problem error")
		}
		err1, err2, err3 := c.SaveUploadedFile(data.Description, filepath.Join(path, service.MakeProblemFileName(problem.ID, problem.Version, "description"))), c.SaveUploadedFile(data.Input, filepath.Join(path, service.MakeProblemFileName(problem.ID, problem.Version, "input"))), c.SaveUploadedFile(data.Output, filepath.Join(path, service.MakeProblemFileName(problem.ID, problem.Version, "output")))
		if err1 != nil || err2 != nil || err3 != nil {
			// 保存文件失败，回滚数据库
			service.SaveProblem(&problemOrigin)
			global.LOG.Warn("save problem " + problem.Name + " file error")
			c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "保存题目文件失败"})
			return
		}
		c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新题目成功"})
	}
}

// DeleteProblem
// @Summary      删除题目
// @Description  删除一个题目
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        id       path      int                true  "题目ID"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/problems/{id} [delete]
func DeleteProblem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if _, notFound := service.GetProblemByID(id); notFound {
		c.JSON(http.StatusOK, model.CommonA{
			Success: false,
			Message: "找不到该题目的信息"})
	} else {
		if err := service.DeleteProblemByID(id); err != nil {
			global.LOG.Panic("DeleteProblem: delete problem error")
		}
		c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除题目成功"})
	}
}

// GetProblemVersion
// @Summary      获取题目版本
// @Description  获取一个题目的版本
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                    true  "token"
// @Param        id       path      int                       true  "题目ID"
// @Success      200      {object}  model.GetProblemVersionA  "是否成功，返回信息，题目版本"
// @Router       /api/v1/problems/{id}/version [get]
func GetProblemVersion(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if problem, notFound := service.GetProblemByID(id); notFound {
		c.JSON(http.StatusOK, model.GetProblemVersionA{
			Success: false,
			Message: "找不到该题目的信息"})
	} else {
		c.JSON(http.StatusOK, model.GetProblemVersionA{Success: true, Message: "获取题目版本成功", Version: problem.Version})
	}
}

// GetProblemList
// @Summary      获取题目列表
// @Description  获取用户所能查看的题目列表
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                 true  "token"
// @Param        page     query     int                    true  "用户位于哪一页，页数从1开始"
// @Param        keyWord  query     string                 true  "当前的(题目名称)搜索关键字，为空字符串表示没有关键字，模糊匹配"
// @Param        sorter   query     int                    true  "用户想按什么排序，1为按题号升序，-1为按题号降序，2为按名称升序，-2为按名称降序，3为按难度升序，-3为按难度降序"
// @Success      200      {object}  model.GetProblemListA  "是否成功，返回信息，题目列表"
// @Router       /api/v1/problems [get]
func GetProblemList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Request.FormValue("page"))
	sorter, _ := strconv.Atoi(c.Request.FormValue("sorter"))
	problems := service.GetReadableProblems(c)
	// TODO 题目名称搜索关键字，模糊查找

	pagedProblems := service.GetProblemsByPage(problems, page, sorter)
	c.JSON(http.StatusOK, model.GetProblemListA{Success: true, Message: "获取成功", ProblemList: pagedProblems, Total: len(problems)})
}
