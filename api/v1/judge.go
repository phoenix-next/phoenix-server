package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model/api"
	"github.com/phoenix-next/phoenix-server/service"
	"net/http"
	"path/filepath"
)

// CreateProblem
// @Summary      创建题目
// @Description  创建一个题目，题目需要包含答案和题面
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.CreateProblemQ  true  "题目名称，题目难度，可读权限，可写权限，组织ID，输入文件，输出文件，题目描述"
// @Success      200   {object}  api.CommonA         "是否成功，返回信息"
// @Router       /api/v1/problems [post]
func CreateProblem(c *gin.Context) {
	var data api.CreateProblemQ
	path := filepath.Join("./", "resource", "problems")
	err := c.ShouldBindJSON(&data)
	if err != nil {
		global.LOG.Panic("CreateProblem: bind data error")
	}
	problem, err := service.CreateProblem(&data)
	if err != nil {
		global.LOG.Panic("CreateProblem: create problem error")
	}
	err1, err2, err3 := c.SaveUploadedFile(data.Description, path+service.MakeProblemFileName(problem.ID, 1, "description")), c.SaveUploadedFile(data.Input, path+service.MakeProblemFileName(problem.ID, 1, "input")), c.SaveUploadedFile(data.Output, path+service.MakeProblemFileName(problem.ID, 1, "output"))
	if err1 != nil || err2 != nil || err3 != nil {
		global.LOG.Panic("save problem " + problem.Name + " file error")
		c.JSON(http.StatusInternalServerError, api.CommonA{Success: false, Message: "保存文件出错"})
		return
	}
	c.JSON(http.StatusOK, api.CommonA{Success: true, Message: "创建题目成功"})
}

// GetProblem
// @Summary      下载题目
// @Description  下载一个题目的信息
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        id    path      int              true  "题目ID"
// @Param        data  body      api.GetProblemQ  true  "GetProblemQ"
// @Success      200   {object}  api.GetProblemA  "是否成功，返回信息"
// @Router       /api/v1/problems/{id} [get]
func GetProblem(c *gin.Context) {

}

// UpdateProblem
// @Summary      更新题目
// @Description  更新一个题目的信息，并自动更新题目版本
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        id    path      int                 true  "题目ID"
// @Param        data  body      api.UpdateProblemQ  true  "UpdateProblemQ"
// @Success      200   {object}  api.CommonA         "是否成功，返回信息"
// @Router       /api/v1/problems/{id} [put]
func UpdateProblem(c *gin.Context) {

}

// DeleteProblem
// @Summary      删除题目
// @Description  删除一个题目
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        id    path      int                 true  "题目ID"
// @Param        data  body      api.DeleteProblemQ  true  "DeleteProblemQ"
// @Success      200   {object}  api.CommonA         "是否成功，返回信息"
// @Router       /api/v1/problems/{id} [delete]
func DeleteProblem(c *gin.Context) {

}

// GetProblemVersion
// @Summary      获取题目版本
// @Description  获取一个题目的版本
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        id    path      int                     true  "题目ID"
// @Param        data  body      api.GetProblemVersionQ  true  "GetProblemVersionQ"
// @Success      200   {object}  api.GetProblemVersionA  "是否成功，返回信息"
// @Router       /api/v1/problems/{id}/version [get]
func GetProblemVersion(c *gin.Context) {

}

// GetProblemList
// @Summary      获取题目列表
// @Description  获取用户所能查看的题目列表
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.GetProblemListQ  true  "GetProblemListQ"
// @Success      200   {object}  api.GetProblemListA  "是否成功，返回信息"
// @Router       /api/v1/problems [get]
func GetProblemList(c *gin.Context) {

}
