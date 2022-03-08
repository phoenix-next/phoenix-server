package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/model"
	"net/http"
)

// CreateContest
// @Summary      创建比赛
// @Description  组织管理员创建一个比赛，比赛包含指定题号的题目
// @Tags         比赛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        data     body      model.CreateContestQ  true  "组织ID，比赛名称，比赛简介，可读权限，比赛包含的题目ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/contests [post]
func CreateContest(c *gin.Context) {
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "TODO"})
}

// GetContest
// @Summary      获取比赛信息
// @Description  获取一个比赛的详细信息，包括该比赛的名称以及包含题目等信息
// @Tags         比赛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       path      int            true  "比赛ID"
// @Success      200      {object}  model.GetContestA  "是否成功，返回信息，比赛ID，比赛名称，比赛简介，比赛包含的题目ID"
// @Router       /api/v1/contests/{id} [get]
func GetContest(c *gin.Context) {
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "TODO"})
}

// UpdateContest
// @Summary      更新比赛
// @Description  管理员更新一个比赛的信息，如更新比赛包含的题目、比赛名称等信息
// @Tags         比赛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        data     body      model.UpdateContestQ  true  "比赛名称，比赛简介，比赛包含的题目ID"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/contests/{id} [put]
func UpdateContest(c *gin.Context) {
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "TODO"})
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
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "TODO"})
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
// @Success      200      {object}  model.GetContestListA  "是否成功，返回信息，题目列表"
// @Router       /api/v1/contests [get]
func GetContestList(c *gin.Context) {
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "TODO"})
}
