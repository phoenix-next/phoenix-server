package database

import "time"

// 帖子
type Post struct {
	ID       uint64 `gorm:"primary_key; not null;" json:"id"`
	UserID   uint64 `gorm:" not null;" json:"user_id"`
	Username string `gorm:"type:varchar(32)" json:"username"`
}

// 评论
type Comment struct {
	ID          uint64    `gorm:"primary_key; not null;" json:"id"`
	ToID        uint64    `gorm:"size:32;" json:"to_id"` // 被评论的id，可为空
	Content     string    `gorm:"size:32;" json:"content"`
	CommentTime time.Time `gorm:"type:datetime" json:"comment_time"`
	UserID      uint64    `gorm:" not null;" json:"user_id"`
	Username    string    `gorm:"type:varchar(32)" json:"username"`
}
