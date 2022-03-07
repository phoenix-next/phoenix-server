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

// CreateTutorial
// @Summary      创建教程
// @Description  创建一个教程，一个教程目前等价于一个markdown文件
// @Tags         教程模块
// @Accept       multipart/form-data
// @Produce      json
// @Param        x-token  header    string                  true  "token"
// @Param        data     body      model.CreateTutorialQ  true  "组织ID，教程名称，教程简介，可读权限，可写权限，教程文件"
// @Success      200      {object}  model.CommonA          "是否成功，返回信息"
// @Router       /api/v1/tutorials [post]
func CreateTutorial(c *gin.Context) {
	// 获取数据
	data := utils.BindJsonData(c, &model.CreateTutorialQ{}).(*model.CreateTutorialQ)
	user := utils.SolveUser(c)
	tutorial := model.Tutorial{Name: data.Name, OrgID: data.OrgID, CreatorID: user.ID, CreatorName: user.Name, Profile: data.Profile, Version: 1, Readable: data.Readable, Writable: data.Writable}
	if err := service.SaveTutorial(&tutorial); err != nil {
		global.LOG.Warn("CreateTutoria;: create tutorial error")
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "创建教程失败"})
		return
	}
	if err := c.SaveUploadedFile(data.File, filepath.Join(global.VP.GetString("root_path"), "resource", "tutorials", service.GetTutorialFileName(tutorial))); err != nil {
		//TODO 发生错误，回滚数据库
		global.LOG.Panic("CreateProblem: save problem error")
	}
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "创建教程成功"})
}

// GetTutorial
// @Summary      下载教程
// @Description  下载一个教程的信息
// @Tags         教程模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string              true  "token"
// @Param        id       path      int                    true  "教程ID"
// @Success      200      {object}  model.GetTutorialA  "组织ID，创建者ID，创建者名称，教程名称，教程简介，教程版本，教程下载路径"
// @Router       /api/v1/tutorials/{id} [get]
func GetTutorial(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if tutorial, notFound := service.GetTutorialByID(id); notFound {
		c.JSON(http.StatusOK, model.GetTutorialA{Success: false, Message: "找不到该教程的信息"})
	} else {
		c.JSON(http.StatusOK, model.GetTutorialA{
			Success:      true,
			Message:      "查找教程成功",
			OrgID:        tutorial.OrgID,
			CreatorID:    tutorial.CreatorID,
			CreatorName:  tutorial.Name,
			Name:         tutorial.Name,
			Profile:      tutorial.Profile,
			Version:      tutorial.Version,
			TutorialPath: filepath.Join(global.VP.GetString("root_path"), "resource", "tutorials", service.GetTutorialFileName(tutorial)),
		})
	}
}

// UpdateTutorial
// @Summary      更新教程
// @Description  更新一个教程的信息，并自动更新教程版本
// @Tags         教程模块
// @Accept       multipart/form-data
// @Produce      json
// @Param        x-token  header    string                 true  "token"
// @Param        id       path      int                 true  "教程ID"
// @Param        data     body      model.UpdateTutorialQ  true  "教程名称，教程简介，可读权限，可写权限，教程文件"
// @Success      200      {object}  model.CommonA          "是否成功，返回信息"
// @Router       /api/v1/tutorials/{id} [put]
func UpdateTutorial(c *gin.Context) {
	// TODO: 逻辑补全
	c.JSON(http.StatusOK, gin.H{"TODO": "remaining logic"})
}

// DeleteTutorial
// @Summary      删除教程
// @Description  删除一个教程
// @Tags         教程模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int            true  "教程ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/tutorials/{id} [delete]
func DeleteTutorial(c *gin.Context) {
	// TODO: 逻辑补全
	c.JSON(http.StatusOK, gin.H{"TODO": "remaining logic"})
}

// GetTutorialVersion
// @Summary      获取教程版本
// @Description  获取一个教程的版本
// @Tags         教程模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                     true  "token"
// @Param        id       path      int                        true  "教程ID"
// @Success      200      {object}  model.GetTutorialVersionA  "是否成功，返回信息，教程版本"
// @Router       /api/v1/tutorials/{id}/version [get]
func GetTutorialVersion(c *gin.Context) {
	// TODO: 逻辑补全
	c.JSON(http.StatusOK, gin.H{"TODO": "remaining logic"})
}

// GetTutorialList
// @Summary      获取教程列表
// @Description  获取用户所能查看的教程列表
// @Tags         教程模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                 true  "token"
// @Param        page     query     int                     true  "用户位于哪一页，页数从1开始"
// @Param        keyWord  query     string                  true  "当前的(教程名称)搜索关键字，为空字符串表示没有关键字，模糊匹配"
// @Param        sorter   query     int                     true  "用户想按什么排序，1为按ID升序，-1为按ID降序，2为按名称升序，-2为按名称降序"
// @Success      200      {object}  model.GetTutorialListA  "是否成功，返回信息，教程列表"
// @Router       /api/v1/tutorials [get]
func GetTutorialList(c *gin.Context) {
	// TODO: 逻辑补全
	c.JSON(http.StatusOK, gin.H{"TODO": "remaining logic"})
}
