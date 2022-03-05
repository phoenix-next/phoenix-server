package database

import "time"

// Post 帖子
type Post struct {
	ID          uint64    `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	OrgID       uint64    `gorm:"not null;" json:"orgID"`
	CreatorID   uint64    `gorm:"not null;" json:"creatorID"`
	Title       string    `gorm:"not null;" json:"title"`
	Content     string    `gorm:"not null;" json:"content"`
	CreatedTime time.Time `gorm:"autoCreateTime" json:"createdTime"`
}

// Comment 评论
type Comment struct {
	ID          uint64    `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	ToID        uint64    `json:"toID"`                    // 被评论的id，可为空
	UserID      uint64    `gorm:"not null;" json:"userID"` // 评论者ID
	Content     string    `gorm:"not null;" json:"content"`
	CommentTime time.Time `gorm:"autoCreateTime" json:"commentTime"`
}
