package database

import "time"

// Post 帖子
type Post struct {
	ID          uint64    `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	OrgID       uint64    `gorm:"not null;" json:"orgID"`
	CreatorID   uint64    `gorm:"not null;" json:"creatorID"`
	CreatorName string    `gorm:"not null;" json:"creatorName"`
	Type        int       `gorm:"not null;" json:"type"` // 帖子所属的板块，0为公告板块，1为划水板块，2为讨论板块
	Title       string    `gorm:"not null;" json:"title"`
	Content     string    `gorm:"not null;" json:"content"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime; not null;" json:"updatedAt"`
}

// Comment 评论
type Comment struct {
	ID          uint64    `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	ToID        uint64    `json:"toID"`                       // 被评论的评论ID，可为空
	PostID      uint64    `json:"postID"`                     // 帖子的ID，即该评论位于哪个帖子下
	CreatorID   uint64    `gorm:"not null;" json:"creatorID"` // 评论者ID
	CreatorName string    `gorm:"not null;" json:"creatorName"`
	Content     string    `gorm:"not null;" json:"content"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime; not null;" json:"updatedAt"`
}
