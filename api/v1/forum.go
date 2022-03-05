package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"github.com/phoenix-next/phoenix-server/utils"
	"net/http"
)

// CreatePost
// @Summary      新建帖子
// @Description  新建一个帖子，需要选择帖子所在的组织，所在的版块
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        data     body      model.CreatePostQ  true  "组织ID，帖子所属板块，帖子标题，帖子内容"
// @Success      200      {object}  model.CommonA      "是否成功，返回信息"
// @Router       /api/v1/posts [post]
func CreatePost(c *gin.Context) {
	data := utils.BindJsonData(c, &model.CreatePostQ{}).(*model.CreatePostQ)
	user := utils.SolveUser(c)
	post := model.Post{Content: data.Content, OrgID: data.OrgID, CreatorID: user.ID, CreatorName: user.Name, Type: data.Type, Title: data.Title}
	if err := global.DB.Create(&post).Error; err != nil {
		global.LOG.Panic("CreatePost: can create post")
	}
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "发帖成功"})
}

// DeletePost
// @Summary      删除帖子
// @Description  删除一个帖子，删除者可以是帖子创建者或者组织管理员
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int            true  "帖子ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/posts/{id} [delete]
func DeletePost(c *gin.Context) {
	//user := utils.SolveUser(c)
	//var post model.Post
	//id, _ := strconv.Atoi(c.Param("id"))
	//global.DB.First(&post, id).Error
	//global.DB.Delete(&post)
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删帖成功"})
}

// UpdatePost
// @Summary      更新帖子内容
// @Description  更新一个帖子的内容
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       path      int                true  "帖子ID"
// @Param        data     body      model.UpdatePostQ  true  "帖子所属板块，帖子标题，帖子内容"
// @Success      200      {object}  model.CommonA      "是否成功，返回信息"
// @Router       /api/v1/posts/{id} [put]
func UpdatePost(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// GetPost
// @Summary      获取帖子详细信息
// @Description  获取一个帖子的详细信息
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string          true  "token"
// @Param        id       path      int             true  "帖子ID"
// @Success      200      {object}  model.GetPostA  "创建者ID，创建者名字，创建者头像路径，标题，内容，最后更新时间"
// @Router       /api/v1/posts/{id} [get]
func GetPost(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// GetAllPost
// @Summary      获取所有帖子
// @Description  给定帖子板块和帖子所属组织，获取该组织该论坛板块的帖子
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       query     int                true  "组织ID"
// @Param        type     query     int                true  "帖子板块"
// @Param        page     query     int                true  "位于第几页，页数从1开始"
// @Success      200      {object}  model.GetAllPostA  "是否成功，返回信息"
// @Router       /api/v1/posts [get]
func GetAllPost(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// CreateComment
// @Summary      新建评论
// @Description  在一个帖子下新建评论
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        id       path      int                   true  "帖子ID"
// @Param        data     body      model.CreateCommentQ  true  "回复的评论ID，用户ID，评论内容"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/posts/{id}/comments [post]
func CreateComment(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// UpdateComment
// @Summary      更新评论
// @Description  更新一条评论的内容
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        id       path      int                   true  "评论ID"
// @Param        data     body      model.UpdateCommentQ  true  "评论内容"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/comments/{id} [put]
func UpdateComment(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// DeleteComment
// @Summary      删除评论
// @Description  删除一条评论
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int            true  "评论ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}

// GetComment
// @Summary      获取评论
// @Description  获取一个帖子下的所有评论
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       path      int                true  "帖子ID"
// @Success      200      {object}  model.GetCommentA  "是否成功，返回信息，评论列表"
// @Router       /api/v1/posts/{id}/comments [get]
func GetComment(c *gin.Context) {
	// TODO 逻辑实现
	c.JSON(http.StatusOK, c.GetString("organization"))
}
