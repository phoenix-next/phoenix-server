package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"github.com/phoenix-next/phoenix-server/service"
	"github.com/phoenix-next/phoenix-server/utils"
	"math"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
		global.LOG.Warn("CreateTutorial;: create tutorial error")
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "创建教程失败"})
		return
	}
	if err := c.SaveUploadedFile(data.File, filepath.Join(global.VP.GetString("root_path"), "resource", "tutorials", service.GetTutorialFileName(tutorial))); err != nil {
		// 回滚数据库
		_ = service.DeleteTutorialByID(tutorial.ID)
		global.LOG.Panic("CreateTutorial: save tutorial error")
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
	// 获取数据
	data := utils.BindJsonData(c, &model.UpdateTutorialQ{}).(*model.UpdateTutorialQ)
	//user := utils.SolveUser(c)
	// TODO 判断可写权限
	if tutorial, notFound := service.GetTutorialByID(data.ID); notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "找不到该教程的信息"})
	} else {
		tutorialOrigin := tutorial
		if err := service.UpdateTutorial(&tutorial, data); err != nil {
			global.LOG.Panic("UpdateTutorial: save tutorial error")
		}
		if err := c.SaveUploadedFile(data.File, filepath.Join(global.VP.GetString("root_path"), "resource", "tutorials", service.GetTutorialFileName(tutorial))); err != nil {
			// 保存文件失败，回滚数据库
			_ = service.SaveTutorial(&tutorialOrigin)
			global.LOG.Warn("save tutorial " + tutorialOrigin.Name + " file error")
			c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "保存教程文件失败"})
			return
		}
		c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新教程成功"})
	}
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
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if _, notFound := service.GetTutorialByID(id); notFound {
		c.JSON(http.StatusOK, model.CommonA{
			Success: false,
			Message: "找不到该教程的信息"})
	} else {
		if err := service.DeleteTutorialByID(id); err != nil {
			global.LOG.Panic("DeleteTutorial: delete tutorial error")
		}
		c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除教程成功"})
	}
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
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if tutorial, notFound := service.GetTutorialByID(id); notFound {
		c.JSON(http.StatusOK, model.GetTutorialVersionA{
			Success: false,
			Message: "找不到该教程的信息"})
	} else {
		c.JSON(http.StatusOK, model.GetProblemVersionA{Success: true, Message: "获取题目版本成功", Version: tutorial.Version})
	}
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
	allTutorials := service.GetAllTutorials()
	page, _ := strconv.Atoi(c.Request.FormValue("page"))
	sorter, _ := strconv.Atoi(c.Request.FormValue("sorter"))
	keyWord := c.Request.FormValue("keyWord")
	// TODO 教程名称搜索关键字，模糊查找
	fuzzyTutorials := make([]model.Tutorial, 0)
	for _, tutorial := range allTutorials {
		if fuzzy.Match(keyWord, tutorial.Name) {
			fuzzyTutorials = append(fuzzyTutorials, tutorial)
		}
	}
	size := 10
	sort.Slice(fuzzyTutorials, func(i, j int) bool {
		if math.Abs(float64(sorter)) == 1 {
			return fuzzyTutorials[i].ID > fuzzyTutorials[j].ID && sorter > 0
		} else if math.Abs(float64(sorter)) == 2 {
			return strings.Compare(fuzzyTutorials[i].Name, fuzzyTutorials[j].Name) < 0 && sorter > 0
		} else {
			return true
		}
	})
	tutorials := fuzzyTutorials[(page-1)*size : int(math.Min(float64(page*size), float64(len(fuzzyTutorials))))]

	c.JSON(http.StatusOK, model.GetTutorialListA{Success: true, Message: "获取成功", TutorialList: tutorials, Total: len(tutorials)})
}
