package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
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
// @Param        x-token  header    string                      true  "token"
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
		Name:       data.Name,
		Version:    1,
		Difficulty: data.Difficulty,
		Readable:   data.Readable,
		Writable:   data.Writable,
		OrgID:      data.OrgID,
		Creator:    user.ID}
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
		_ = service.DeleteProblemByID(problem.ID)
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
// @Success      200      {object}  model.GetProblemA  "题目名称，题目难度，输入文件，输出文件，题目描述，评测结果"
// @Router       /api/v1/problems/{id} [get]
func GetProblem(c *gin.Context) {
	// 获取请求参数
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.GetProblemA{Success: false, Message: "请求参数非法"})
		return
	}
	// 题目的存在性判定
	problem, notFound := service.GetProblemByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.GetProblemA{Success: false, Message: "找不到该题目的信息"})
		return
	}
	// 用户权限判定
	if !service.JudgeReadPermission(problem.OrgID, problem.Readable, problem.Creator, c) {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "您对该题目无可读权限"})
		return
	}
	// 返回结果
	c.JSON(http.StatusOK, model.GetProblemA{
		Success:     true,
		Message:     "",
		Name:        problem.Name,
		Difficulty:  problem.Difficulty,
		Input:       service.GetProblemFileUrl(&problem, "input"),
		Output:      service.GetProblemFileUrl(&problem, "output"),
		Description: service.GetProblemFileUrl(&problem, "description"),
		Result:      service.GetUserFinalJudge(utils.SolveUser(c).ID, id),
	})

}

// UpdateProblem
// @Summary      更新题目
// @Description  更新一个题目的信息，并自动更新题目版本
// @Tags         评测模块
// @Accept       multipart/form-data
// @Produce      json
// @Param        x-token  header    string                   true  "token"
// @Param        data     body      model.UpdateProblemQ  true  "题目ID，题目名称，题目难度，可读权限，可写权限，组织ID，输入文件，输出文件，题目描述"
// @Success      200      {object}  model.CommonA               "是否成功，返回信息"
// @Router       /api/v1/problems/{id} [put]
func UpdateProblem(c *gin.Context) {
	// 获取请求数据
	var data model.UpdateProblemQ
	err := c.ShouldBind(&data)
	if err != nil {
		global.LOG.Panic("UpdateProblem: bind data error")
	}
	// 题目不存在的情况
	problem, notFound := service.GetProblemByID(data.ID)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "找不到该题目的信息"})
		return
	}
	// 用户没有权限修改题目的情况
	if !service.JudgeWritePermission(problem.OrgID, problem.Writable, problem.Creator, c) {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "您对该题目无可写权限"})
		return
	}
	// 更新题目，同时保存原来的题目以备回滚之需
	problemOrigin := problem
	err = service.UpdateProblem(&problem, &data)
	if err != nil {
		global.LOG.Panic("UpdateProblem: save problem error")
	}
	// 保存题目相关的文件
	dir := global.VP.GetString("problem_path")
	err1 := c.SaveUploadedFile(data.Description, filepath.Join(dir, service.MakeProblemFileName(problem.ID, problem.Version, "description")))
	err2 := c.SaveUploadedFile(data.Input, filepath.Join(dir, service.MakeProblemFileName(problem.ID, problem.Version, "input")))
	err3 := c.SaveUploadedFile(data.Output, filepath.Join(dir, service.MakeProblemFileName(problem.ID, problem.Version, "output")))
	// 保存文件失败，回滚数据库
	if err1 != nil || err2 != nil || err3 != nil {
		service.SaveProblem(&problemOrigin)
		global.LOG.Warn("save problem " + problem.Name + " file error")
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "保存题目文件失败"})
		return
	}
	// 成功更新题目
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新题目成功"})
}

// DeleteProblem
// @Summary      删除题目
// @Description  删除一个题目
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int                true  "题目ID"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/problems/{id} [delete]
func DeleteProblem(c *gin.Context) {
	// 获取请求参数
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "请求参数非法"})
		return
	}
	// 题目的存在性判定
	problem, notFound := service.GetProblemByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "找不到该题目的信息"})
		return
	}
	// 删除题目
	if err = global.DB.Delete(&problem).Error; err != nil {
		global.LOG.Panic("DeleteProblem: delete problem error")
	}
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除题目成功"})

}

// GetProblemVersion
// @Summary      获取题目版本
// @Description  获取一个题目的版本
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                    true  "token"
// @Param        id       path      int                      true  "题目ID"
// @Success      200      {object}  model.GetProblemVersionA  "是否成功，返回信息，题目版本"
// @Router       /api/v1/problems/{id}/version [get]
func GetProblemVersion(c *gin.Context) {
	// 获取请求数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.GetProblemVersionA{Success: false, Message: "找不到该题目的信息"})
		return
	}
	// 题目的存在性判定
	problem, notFound := service.GetProblemByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.GetProblemVersionA{Success: false, Message: "找不到该题目的信息"})
		return
	}
	c.JSON(http.StatusOK, model.GetProblemVersionA{Success: true, Message: "", Version: problem.Version})

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
	// 获取请求数据
	user := utils.SolveUser(c)
	page, err1 := strconv.Atoi(c.Query("page"))
	sorter, err2 := strconv.Atoi(c.Query("sorter"))
	keyWord := c.Query("keyWord")
	// 请求数据不合法的情况
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, model.GetProblemListA{Success: false, Message: "请求参数不合法"})
		return
	}
	// 获取可读的题目
	problems := service.GetReadableProblems(c)
	// 对题目标题进行模糊查找
	var resProblems []model.Problem
	for _, problem := range problems {
		if fuzzy.Match(keyWord, problem.Name) {
			resProblems = append(resProblems, problem)
		}
	}
	// 找不到题目的情况
	if len(resProblems) == 0 {
		c.JSON(http.StatusOK, model.GetProblemListA{
			Success:     true,
			Message:     "",
			ProblemList: []model.ProblemT{},
			Total:       0})
		return
	}
	// 对题目进行分页
	pagedProblems := service.GetProblemsByPage(resProblems, page, sorter)
	// 获取用户的评测结果
	var finalProblems []model.ProblemT
	for _, problem := range pagedProblems {
		finalProblems = append(finalProblems, model.ProblemT{
			ProblemID:   problem.ID,
			ProblemName: problem.Name,
			Difficulty:  problem.Difficulty,
			Result:      service.GetUserFinalJudge(user.ID, problem.ID)})
	}
	// 返回响应
	c.JSON(http.StatusOK, model.GetProblemListA{
		Success:     true,
		Message:     "",
		ProblemList: finalProblems,
		Total:       len(resProblems)})
}

// UploadProblemRecord
// @Summary      上传评测结果
// @Description  上传一个题目的评测结果，用户必须有该题目的读权限
// @Tags         评测模块
// @Accept       multipart/form-data
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        id       path      int                         true  "题目ID"
// @Param        data     body      model.UploadProblemRecordQ  true  "评测结果，代码文件"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/problems/{id}/records [post]
func UploadProblemRecord(c *gin.Context) {
	// 获取请求数据
	var data model.UploadProblemRecordQ
	if err := c.ShouldBind(&data); err != nil {
		global.LOG.Panic("UploadProblemRecord: bind data error")
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.GetProblemVersionA{Success: false, Message: "题目ID不为整数"})
		return
	}
	// 题目的存在性判定
	problem, notFound := service.GetProblemByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.GetProblemA{Success: false, Message: "找不到该题目的信息"})
		return
	}
	user := utils.SolveUser(c)
	// 用户权限判定
	if !service.JudgeReadPermission(problem.OrgID, problem.Readable, problem.Creator, c) {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "您对该题目无可读权限"})
		return
	}
	result := model.Result{Result: data.Result, UserID: user.ID, ProblemID: problem.ID}
	if err = global.DB.Create(&result).Error; err != nil {
		global.LOG.Warn("UploadProblemRecord: judge problem error")
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "评测题目失败"})
		return
	}
	if err := c.SaveUploadedFile(data.Code, filepath.Join(global.VP.GetString("code_path"), service.GetCodeFileName(result))); err != nil {
		//发生错误，回滚数据库
		_ = global.DB.Where("id = ?", result.ID).Delete(model.Result{}).Error
		global.LOG.Panic("UploadProblemRecord: save judge code error")
	}
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "评测题目成功"})
}

// GetProblemRecord
// @Summary      获取评测结果
// @Description  获取一个题目的评测结果，用户必须有该题目的读权限
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        id       path      int                       true  "题目ID"
// @Success      200      {object}  model.GetProblemRecordA  "是否成功，返回信息"
// @Router       /api/v1/problems/{id}/records [get]
func GetProblemRecord(c *gin.Context) {

}
