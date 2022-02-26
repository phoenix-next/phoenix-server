package v1

import "github.com/gin-gonic/gin"

// CreateProblem
// @Summary      创建题目
// @Description  创建一个题目，题目需要包含答案和题面
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.CreateProblemQ  true  ""
// @Success      200   {object}  api.CommonA    "是否成功，返回信息"
// @Router       /api/v1/problems [post]
func CreateProblem(c *gin.Context) {

}

// GetProblem
// @Summary      下载题目
// @Description  下载一个题目的信息
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.GetProblemQ  true  ""
// @Success      200   {object}  api.CommonA    "是否成功，返回信息"
// @Router       /api/v1/problems/:id [get]
func GetProblem(c *gin.Context) {

}

// UpdateProblem
// @Summary      更新题目
// @Description  更新一个题目的信息，并自动更新题目版本
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.UpdateProblemQ  true  ""
// @Success      200   {object}  api.CommonA    "是否成功，返回信息"
// @Router       /api/v1/problems/:id [put]
func UpdateProblem(c *gin.Context) {

}

// DeleteProblem
// @Summary      删除题目
// @Description  删除一个题目
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.DeleteProblemQ  true  ""
// @Success      200   {object}  api.CommonA    "是否成功，返回信息"
// @Router       /api/v1/problems/:id [delete]
func DeleteProblem(c *gin.Context) {

}

// GetProblemVersion
// @Summary      获取题目版本
// @Description  获取一个题目的版本
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.GetProblemVersionQ  true  ""
// @Success      200   {object}  api.CommonA    "是否成功，返回信息"
// @Router       /api/v1/problems/:id/version [get]
func GetProblemVersion(c *gin.Context) {

}

// GetProblemList
// @Summary      获取题目列表
// @Description  获取用户所能查看的题目列表
// @Tags         评测模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.GetProblemListQ  true  ""
// @Success      200   {object}  api.CommonA    "是否成功，返回信息"
// @Router       /api/v1/problems [get]
func GetProblemList(c *gin.Context) {

}
