package api

import "time"

type CreatePostQ struct {
	OrgID   uint64 `json:"orgID"`
	Type    int    `json:"type"` // 帖子所属的板块，0为公告板块，1为划水板块，2为讨论板块
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostQ struct {
	Type    int    `json:"type"` // 帖子所属的板块，0为公告板块，1为划水板块，2为讨论板块
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GetPostA struct {
	CreatorID     uint64    `json:"creatorID"`
	CreatorName   string    `json:"creatorName"`
	CreatorAvatar string    `json:"creatorAvatar"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type GetAllPostA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Total   int    `json:"total"`
	Posts   []struct {
		ID            uint64    `json:"id"`
		CreatorID     uint64    `gorm:"not null;" json:"creatorID"`
		CreatorName   string    `json:"creatorName"`
		CreatorAvatar string    `json:"creatorAvatar"`
		Title         string    `json:"title"`
		UpdatedAt     time.Time `json:"updatedAt"`
	}
}

type CreateCommentQ struct {
	ToID      uint64 `json:"toID"`      // 被评论的评论ID，可为空
	CreatorID uint64 `json:"creatorID"` // 评论者ID
	Content   string `json:"content"`
}

type UpdateCommentQ struct {
	Content string `json:"content"`
}

type GetCommentA struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Comments []struct {
		ID            uint64    `json:"id"`
		ToID          uint64    `json:"toID"`      // 被评论的评论ID，可为空
		CreatorID     uint64    `json:"creatorID"` // 评论者ID
		CreatorName   string    `json:"creatorName"`
		CreatorAvatar string    `json:"creatorAvatar"`
		Content       string    `json:"content"`
		UpdatedAt     time.Time `json:"updatedAt"`
	}
}
