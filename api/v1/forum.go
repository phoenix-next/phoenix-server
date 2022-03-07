package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"github.com/phoenix-next/phoenix-server/service"
	"github.com/phoenix-next/phoenix-server/utils"
	"net/http"
	"strconv"
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
	// 获取数据
	data := utils.BindJsonData(c, &model.CreatePostQ{}).(*model.CreatePostQ)
	user := utils.SolveUser(c)
	// 用户权限判定
	ok, err := service.IsUserInThisOrganization(user.ID, data.OrgID)
	if !ok || err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "用户没有权限在此发帖"})
		return
	}
	// 成功新建帖子
	post := model.Post{
		Content:       data.Content,
		OrgID:         data.OrgID,
		CreatorID:     user.ID,
		CreatorName:   user.Name,
		CreatorAvatar: user.Avatar,
		Type:          data.Type,
		Title:         data.Title}
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
	// 获取数据
	user := utils.SolveUser(c)
	post, ok := service.GetPostFromParam(c)
	// 帖子存在性判定
	if !ok {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "帖子不存在"})
		return
	}
	// 删帖权限判定
	if post.CreatorID == user.ID {
		global.DB.Delete(&post)
		c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删帖成功"})
		return
	}
	for _, adminID := range service.GetOrganizationAdmin(post.OrgID) {
		if adminID == user.ID {
			global.DB.Delete(&post)
			c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删帖成功"})
			return
		}
	}
	// 删帖失败
	c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "没有权限删帖"})
}

// UpdatePost
// @Summary      更新帖子内容
// @Description  更新一个帖子的内容，更新者可以是帖子创建者或者组织管理员
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       path      int                true  "帖子ID"
// @Param        data     body      model.UpdatePostQ  true  "帖子所属板块，帖子标题，帖子内容"
// @Success      200      {object}  model.CommonA      "是否成功，返回信息"
// @Router       /api/v1/posts/{id} [put]
func UpdatePost(c *gin.Context) {
	// 获取数据
	user := utils.SolveUser(c)
	data := utils.BindJsonData(c, &model.UpdatePostQ{}).(*model.UpdatePostQ)
	post, ok := service.GetPostFromParam(c)
	// 帖子存在性判定
	if !ok {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "帖子不存在"})
		return
	}
	// 更新权限判定
	if post.CreatorID == user.ID {
		global.DB.Model(&post).Updates(model.Post{Content: data.Content, Title: data.Title, Type: data.Type})
		c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新帖子成功"})
		return
	}
	for _, adminID := range service.GetOrganizationAdmin(post.OrgID) {
		if adminID == user.ID {
			global.DB.Model(&post).Updates(model.Post{Content: data.Content, Title: data.Title, Type: data.Type})
			c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "更新帖子成功"})
			return
		}
	}
	// 更新失败
	c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "没有权限更新帖子"})
}

// GetPost
// @Summary      获取帖子详细信息
// @Description  获取一个帖子的详细信息，用户必须是该帖子所在组织的成员
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string          true  "token"
// @Param        id       path      int             true  "帖子ID"
// @Success      200      {object}  model.GetPostA  "创建者ID，创建者名字，创建者头像路径，标题，内容，最后更新时间"
// @Router       /api/v1/posts/{id} [get]
func GetPost(c *gin.Context) {
	// 帖子存在性判定
	post, ok := service.GetPostFromParam(c)
	if !ok {
		c.JSON(http.StatusOK, model.GetPostA{Success: false, Message: "帖子不存在"})
		return
	}
	// 用户权限判定s
	user := utils.SolveUser(c)
	ok, _ = service.IsUserInThisOrganization(user.ID, post.OrgID)
	if !ok {
		c.JSON(http.StatusOK, model.GetPostA{Success: false, Message: "用户没有查阅权限"})
		return
	}
	// 成功获取信息
	creator, _ := service.GetUserByID(post.CreatorID)
	c.JSON(http.StatusOK, model.GetPostA{Success: true,
		Message:       "获取信息成功",
		CreatorID:     post.CreatorID,
		CreatorName:   post.CreatorName,
		CreatorAvatar: creator.Avatar,
		Title:         post.Title,
		Content:       post.Content,
		UpdatedAt:     post.UpdatedAt})
}

// GetAllPost
// @Summary      获取所有帖子
// @Description  给定帖子板块和帖子所属组织，获取该组织该论坛板块的帖子，用户必须是帖子所在组织的成员
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       query     int                true  "组织ID"
// @Param        type     query     int                true  "帖子板块"
// @Param        page     query     int                true  "位于第几页，页数从1开始"
// @Success      200      {object}  model.GetAllPostA  "是否成功，返回信息，帖子总数，帖子列表"
// @Router       /api/v1/posts [get]
func GetAllPost(c *gin.Context) {
	// 获取路径参数
	oid, err1 := strconv.ParseUint(c.Query("id"), 10, 64)
	postType, err2 := strconv.Atoi(c.Query("type"))
	page, err3 := strconv.Atoi(c.Query("page"))
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, model.GetAllPostA{Success: false, Message: "参数非法"})
		return
	}
	// 判断用户权限
	user := utils.SolveUser(c)
	ok, err := service.IsUserInThisOrganization(user.ID, oid)
	if !ok || err != nil {
		c.JSON(http.StatusOK, model.GetAllPostA{Success: false, Message: "用户没有查看帖子权限"})
		return
	}
	// 得到所有帖子，及所有帖子的总数
	posts := service.GetAllPosts(oid, postType)
	totalPage := len(posts) / 10
	if len(posts)%10 != 0 {
		totalPage += 1
	}
	// 没有帖子的情况
	if totalPage == 0 {
		c.JSON(http.StatusOK, model.GetAllPostA{
			Success: true,
			Message: "获取帖子成功",
			Total:   0,
			Posts:   []model.PostT{}})
		return
	}
	// 页数不合法的情况
	if page <= 0 || page > totalPage {
		c.JSON(http.StatusOK, model.GetAllPostA{Success: false, Message: "页数非法"})
		return
	}
	// 获取端点位置，并对帖子切片
	start, end := (page-1)*10, page*10
	if length := len(posts); end > length {
		end = length
	}
	filteredPosts := posts[start:end]
	// 返回
	c.JSON(http.StatusOK, model.GetAllPostA{
		Success: true,
		Message: "获取帖子成功",
		Total:   len(posts),
		Posts:   filteredPosts})
}

// CreateComment
// @Summary      新建评论
// @Description  在一个帖子下新建评论，用户必须是该帖子所在组织的成员
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        id       path      int                   true  "帖子ID"
// @Param        data     body      model.CreateCommentQ  true  "回复的评论ID(若不回复则为0)，评论内容"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/posts/{id}/comments [post]
func CreateComment(c *gin.Context) {
	// 帖子的存在性判定
	post, ok := service.GetPostFromParam(c)
	if !ok {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "帖子不存在"})
		return
	}
	// 用户权限判定
	user := utils.SolveUser(c)
	ok, err := service.IsUserInThisOrganization(user.ID, post.OrgID)
	if !ok || err != nil {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "用户没有权限进行该操作"})
		return
	}
	// 成功创建评论
	data := utils.BindJsonData(c, &model.CreateCommentQ{}).(*model.CreateCommentQ)
	global.DB.Create(&model.Comment{
		OrgID:         post.OrgID,
		Content:       data.Content,
		ToID:          data.ToID,
		PostID:        post.ID,
		CreatorID:     user.ID,
		CreatorName:   user.Name,
		CreatorAvatar: user.Avatar})
	c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "评论成功"})
}

// UpdateComment
// @Summary      更新评论
// @Description  更新一条评论的内容，用户是评论创建者或者组织管理员
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string                true  "token"
// @Param        id       path      int                   true  "评论ID"
// @Param        data     body      model.UpdateCommentQ  true  "评论内容"
// @Success      200      {object}  model.CommonA         "是否成功，返回信息"
// @Router       /api/v1/comments/{id} [put]
func UpdateComment(c *gin.Context) {
	// 评论的存在性判定
	comment, ok := service.GetCommentFromParam(c)
	if !ok {
		c.JSON(http.StatusOK, model.GetPostA{Success: false, Message: "评论不存在"})
		return
	}
	// 用户权限判定
	user := utils.SolveUser(c)
	data := utils.BindJsonData(c, &model.UpdateCommentQ{}).(*model.UpdateCommentQ)
	if user.ID == comment.CreatorID {
		global.DB.Model(&comment).Updates(model.Comment{Content: data.Content})
		c.JSON(http.StatusOK, model.GetPostA{Success: true, Message: "评论更新成功"})
		return
	}
	for _, adminID := range service.GetOrganizationAdmin(comment.OrgID) {
		if adminID == user.ID {
			global.DB.Model(&comment).Updates(model.Comment{Content: data.Content})
			c.JSON(http.StatusOK, model.GetPostA{Success: true, Message: "评论更新成功"})
			return
		}
	}
	// 更新评论失败
	c.JSON(http.StatusOK, model.GetPostA{Success: false, Message: "没有权限更新评论"})
}

// DeleteComment
// @Summary      删除评论
// @Description  删除一条评论，删者是评论创建者或者组织管理员
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string         true  "token"
// @Param        id       path      int            true  "评论ID"
// @Success      200      {object}  model.CommonA  "是否成功，返回信息"
// @Router       /api/v1/comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	// 获取数据
	user := utils.SolveUser(c)
	comment, ok := service.GetCommentFromParam(c)
	// 评论存在性判定
	if !ok {
		c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "评论不存在"})
		return
	}
	// 删除评论权限判定
	if comment.CreatorID == user.ID {
		global.DB.Delete(&comment)
		c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除评论成功"})
		return
	}
	for _, adminID := range service.GetOrganizationAdmin(comment.OrgID) {
		if adminID == user.ID {
			global.DB.Delete(&comment)
			c.JSON(http.StatusOK, model.CommonA{Success: true, Message: "删除评论成功"})
			return
		}
	}
	// 删除评论失败
	c.JSON(http.StatusOK, model.CommonA{Success: false, Message: "没有权限删除评论"})
}

// GetComment
// @Summary      获取评论
// @Description  获取一个帖子下的所有评论，用户必须属于该帖子所在的组织
// @Tags         论坛模块
// @Accept       json
// @Produce      json
// @Param        x-token  header    string             true  "token"
// @Param        id       path      int                true  "帖子ID"
// @Success      200      {object}  model.GetCommentA  "是否成功，返回信息，评论列表"
// @Router       /api/v1/posts/{id}/comments [get]
func GetComment(c *gin.Context) {
	// 帖子存在性判定
	post, ok := service.GetPostFromParam(c)
	if !ok {
		c.JSON(http.StatusOK, model.GetCommentA{Success: false, Message: "该帖子不存在"})
		return
	}
	// 用户权限判定
	user := utils.SolveUser(c)
	ok, err := service.IsUserInThisOrganization(user.ID, post.OrgID)
	if !ok || err != nil {
		c.JSON(http.StatusOK, model.GetCommentA{Success: false, Message: "用户没有权限获取该帖子评论"})
		return
	}
	// 成功获取评论
	comments := service.GetAllCommentByPostID(post.ID)
	c.JSON(http.StatusOK, model.GetCommentA{Success: true, Message: "获取评论成功", Comments: comments})
}
