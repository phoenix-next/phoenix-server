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
// @Param        file     formData  file                   true  "教程文件"
// @Param        data     body      model.CreateTutorialQ  true  "组织ID，教程名称，教程简介，可读权限，可写权限"
// @Success      200      {object}  model.CommonA          "是否成功，返回信息"
// @Router       /api/v1/tutorials [post]
func CreateTutorial(c *gin.Context) {
	// 获取请求数据
	var data model.CreateTutorialQ
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "请求参数非法"})
		return
	}
	user := utils.SolveUser(c)
	// 新建教程
	tutorial := model.Tutorial{
		Name:      data.Name,
		OrgID:     data.OrgID,
		CreatorID: user.ID,
		Profile:   data.Profile,
		Version:   1,
		Readable:  data.Readable,
		Writable:  data.Writable}
	if err := service.SaveTutorial(&tutorial); err != nil {
		global.LOG.Warn("CreateTutorial;: create tutorial error")
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "创建教程失败"})
		return
	}
	if err := c.SaveUploadedFile(data.File, filepath.Join(global.VP.GetString("tutorial_path"), service.GetTutorialFileName(tutorial))); err != nil {
		// 回滚数据库
		_ = service.DeleteTutorialByID(tutorial.ID)
		global.LOG.Panic(err)
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
	// 获取请求中的数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.GetTutorialA{Success: false, Message: "请求参数非法"})
		return
	}
	// 找不到教程的情况
	tutorial, notFound := service.GetTutorialByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.GetTutorialA{Success: false, Message: "找不到该教程的信息"})
		return
	}
	// 当前用户没有读权限的情况
	if !service.JudgeReadPermission(tutorial.OrgID, tutorial.Readable, tutorial.CreatorID, c) {
		c.JSON(http.StatusOK, model.GetTutorialA{Success: false, Message: "您没有可读权限"})
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, model.GetTutorialA{
		Success:      true,
		OrgID:        tutorial.OrgID,
		CreatorID:    tutorial.CreatorID,
		CreatorName:  tutorial.Name,
		Name:         tutorial.Name,
		Profile:      tutorial.Profile,
		Version:      tutorial.Version,
		TutorialPath: "resource/tutorial/" + service.GetTutorialFileName(tutorial)})
}

// UpdateTutorial
// @Summary      更新教程
// @Description  更新一个教程的信息，并自动更新教程版本
// @Tags         教程模块
// @Accept       multipart/form-data
// @Produce      json
// @Param        x-token  header    string                 true  "token"
// @Param        id       path      int                 true  "教程ID"
// @Param        file     formData  file                   true  "教程文件"
// @Param        data     body      model.UpdateTutorialQ  true  "教程名称，教程简介"
// @Success      200      {object}  model.CommonA          "是否成功，返回信息"
// @Router       /api/v1/tutorials/{id} [put]
func UpdateTutorial(c *gin.Context) {
	// 获取数据
	var data model.UpdateTutorialQ
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "请求参数非法"})
		return
	}
	// 教程不存在的情况
	tutorial, notFound := service.GetTutorialByID(data.ID)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "找不到该教程的信息"})
		return
	}
	// 当前用户没有写权限的情况
	if !service.JudgeWritePermission(tutorial.OrgID, tutorial.Writable, tutorial.CreatorID, c) {
		c.JSON(http.StatusOK, model.GetTutorialA{Success: false, Message: "您无可写权限"})
		return
	}
	// 用户具有写权限，则更新数据库，并保存原教程以备回滚数据库之需
	tutorialOrigin := tutorial
	if err := service.UpdateTutorial(&tutorial, &data); err != nil {
		global.LOG.Panic("UpdateTutorial: save tutorial error")
	}
	// 保存文件失败，回滚数据库
	filePath := filepath.Join(global.VP.GetString("tutorial_path"), service.GetTutorialFileName(tutorial))
	if err := c.SaveUploadedFile(data.File, filePath); err != nil {
		_ = service.SaveTutorial(&tutorialOrigin)
		global.LOG.Warn("save tutorial " + tutorialOrigin.Name + " file error")
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "保存教程文件失败"})
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新成功"})
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
	// 获取请求中的数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "请求参数非法"})
		return
	}
	// 找不到教程的情况
	_, notFound := service.GetTutorialByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "找不到该教程的信息"})
		return
	}
	// TODO: 用户权限判定
	// 成功删除教程
	if err = service.DeleteTutorialByID(id); err != nil {
		global.LOG.Panic("DeleteTutorial: delete tutorial error")
	}
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除教程成功"})
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
	// 获取请求中的数据
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.GetTutorialVersionA{Success: false, Message: "请求参数非法"})
		return
	}
	// 找不到教程的情况
	tutorial, notFound := service.GetTutorialByID(id)
	if notFound {
		c.JSON(http.StatusOK, model.GetTutorialVersionA{Success: false, Message: "找不到该教程的信息"})
	}
	// 成功返回响应
	c.JSON(http.StatusOK, model.GetProblemVersionA{Success: true, Version: tutorial.Version})
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
	// 获取请求中的数据
	page, _ := strconv.Atoi(c.Query("page"))
	sorter, _ := strconv.Atoi(c.Query("sorter"))
	keyWord := c.Query("keyWord")
	// 获取所有的教程
	allTutorials := service.GetAllTutorials()
	// 教程名称搜索关键字，模糊查找
	fuzzyTutorials := make([]model.Tutorial, 0)
	for _, tutorial := range allTutorials {
		if fuzzy.MatchFold(keyWord, tutorial.Name) {
			fuzzyTutorials = append(fuzzyTutorials, tutorial)
		}
	}
	// 找不到教程的情况
	if len(fuzzyTutorials) == 0 {
		c.JSON(http.StatusOK, model.GetTutorialListA{
			Success:      true,
			Total:        0,
			TutorialList: make([]model.TutorialT, 0)})
		return
	}
	// 对教程进行排序与分页
	size := 10
	sort.Slice(fuzzyTutorials, func(i, j int) (res bool) {
		switch sorter {
		case 1, -1:
			res = fuzzyTutorials[i].ID > fuzzyTutorials[j].ID
		case 2, -2:
			res = strings.Compare(fuzzyTutorials[i].Name, fuzzyTutorials[j].Name) < 0
		}
		if sorter < 0 {
			res = !res
		}
		return
	})
	tutorials := fuzzyTutorials[(page-1)*size : int(math.Min(float64(page*size), float64(len(fuzzyTutorials))))]
	// 获取创建者名称
	finalTutorials := make([]model.TutorialT, 0)
	for _, tutorial := range tutorials {
		tmp, _ := service.GetUserByID(tutorial.CreatorID)
		finalTutorials = append(finalTutorials, model.TutorialT{
			ID:          tutorial.ID,
			Profile:     tutorial.Profile,
			Name:        tutorial.Name,
			CreatorName: tmp.Name})
	}
	// 返回响应
	c.JSON(http.StatusOK, model.GetTutorialListA{Success: true, TutorialList: finalTutorials, Total: len(tutorials)})
}
