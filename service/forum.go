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

// 数据库操作

// GetPostByID 通过ID获取帖子
func GetPostByID(id uint64) (post model.Post, notFound bool) {
	if err := global.DB.First(&post, id).Error; err != nil {
		return post, true
	}
	return post, false
}

// GetAllPosts 已知组织ID和帖子板块，获取所有帖子
func GetAllPosts(oid uint64, postType int) (posts []model.PostT) {
	global.DB.Model(&model.Post{}).Where("OrgID = ? AND Type = ?", oid, postType).Find(&posts)
	return
}
