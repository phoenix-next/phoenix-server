package service

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"strconv"
)

// Helper

// GetPostFromParam 通过路径参数获取对应的帖子
func GetPostFromParam(c *gin.Context) (post model.Post, ok bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		global.LOG.Warn("GetPostFromParam: id invalid")
		return post, false
	}
	post, notFound := GetPostByID(id)
	if notFound {
		global.LOG.Warn("GetPostFromParam: post not found")
		return post, false
	}
	return post, true
}

// GetCommentFromParam 通过路径参数获取对应的评论
func GetCommentFromParam(c *gin.Context) (comment model.Comment, ok bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		global.LOG.Warn("GetCommentFromParam: id invalid")
		return comment, false
	}
	comment, notFound := GetCommentByID(id)
	if notFound {
		global.LOG.Warn("GetCommentFromParam: comment not found")
		return comment, false
	}
	return comment, true
}

// 数据库操作

// GetPostByID 通过ID获取帖子
func GetPostByID(id uint64) (post model.Post, notFound bool) {
	if err := global.DB.First(&post, id).Error; err != nil {
		return post, true
	}
	return post, false
}

// GetAllPosts 已知组织ID和帖子板块，获取所有帖子
func GetAllPosts(oid uint64, postType int) (posts []model.Post) {
	global.DB.Where("org_id = ? AND type = ?", oid, postType).Find(&posts)
	return
}

// GetCommentByID 通过ID获取评论
func GetCommentByID(id uint64) (comment model.Comment, notFound bool) {
	if err := global.DB.First(&comment, id).Error; err != nil {
		return comment, true
	}
	return comment, false
}

// GetAllCommentByPostID 获取一个帖子下的所有评论
func GetAllCommentByPostID(postID uint64) (comments []model.Comment) {
	global.DB.Where("post_id = ?", postID).Find(&comments)
	return
}
