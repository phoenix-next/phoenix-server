package model

import "time"

// 论坛模块

// Post 帖子
type Post struct {
	ID        uint64    `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	OrgID     uint64    `gorm:"not null;" json:"orgID"`
	CreatorID uint64    `gorm:"not null;" json:"creatorID"`
	Type      int       `gorm:"not null;" json:"type"` // 帖子所属的板块，0为公告板块，1为划水板块，2为讨论板块
	Title     string    `gorm:"not null;" json:"title"`
	Content   string    `gorm:"not null;" json:"content"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; not null;" json:"updatedAt"`
}

// Comment 评论
type Comment struct {
	ID            uint64    `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	OrgID         uint64    `gorm:"not null;" json:"orgID"`
	ToID          uint64    `json:"toID"`                       // 被评论的评论ID，可为空
	PostID        uint64    `json:"postID"`                     // 帖子的ID，即该评论位于哪个帖子下
	CreatorID     uint64    `gorm:"not null;" json:"creatorID"` // 评论者ID
	CreatorName   string    `gorm:"not null;" json:"creatorName"`
	CreatorAvatar string    `json:"creatorAvatar"`
	Content       string    `gorm:"not null;" json:"content"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime; not null;" json:"updatedAt"`
}

// 评测模块

// Contest 比赛
type Contest struct {
	ID        uint64    `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	OrgID     uint64    `gorm:"not null;" json:"orgID"` // 比赛所属的组织ID
	Name      string    `gorm:"size:32; not null;" json:"name"`
	Profile   string    `gorm:"not null;" json:"profile"`
	Readable  int       `gorm:"not null" json:"readable"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// Problem 题目
type Problem struct {
	ID          uint64    `gorm:"primary_key;autoIncrement;not null;" json:"id"`
	OrgID       uint64    `json:"organization"`
	Name        string    `gorm:"size:32; not null" json:"name"`
	Version     int       `gorm:"not null;" json:"version"`
	Difficulty  int       `gorm:"not null" json:"difficulty"`
	Readable    int       `gorm:"not null" json:"readable"`
	Writable    int       `gorm:"not null" json:"writable"`
	Creator     uint64    `gorm:"not null" json:"creator"`
	CreatedTime time.Time `gorm:"autoCreateTime" json:"createdTime"`
}

// 社交模块

// User 用户
type User struct {
	ID       uint64    `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	Name     string    `gorm:"size:32; not null;" json:"name"`
	Email    string    `gorm:"size:32; not null; unique;" json:"email"`
	Password string    `gorm:"size:128; not null;" json:"password"`
	Avatar   string    `json:"avatar"`
	Profile  string    `gorm:"size:256;" json:"profile"`
	RegTime  time.Time `gorm:"autoCreateTime" json:"regTime"`
}

// Captcha 邮箱验证码
type Captcha struct {
	ID       uint64    `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	Email    string    `gorm:"size:32; not null; unique;" json:"email"`
	Type     int       `gorm:"not null;" json:"type"` // 1为注册验证码，2为忘记密码验证码
	Captcha  uint64    `gorm:"not null;" json:"captcha"`
	SendTime time.Time `gorm:"autoCreateTime" json:"sendTime"`
}

// Organization 组织
type Organization struct {
	ID          uint64    `gorm:"primary_key; autoIncrement;not null;" json:"id"`
	Name        string    `gorm:"size:32; not null; unique" json:"name"`
	Profile     string    `gorm:"size:255;" json:"profile"`
	CreatorID   uint64    `gorm:"not null;" json:"creatorID"` // 创建者ID
	CreatorName string    `gorm:"size:32;not null;" json:"creatorName"`
	CreatedTime time.Time `gorm:"autoCreateTime;" json:"createdTime"`
}

// 教程模块

// Tutorial 教程
type Tutorial struct {
	ID          uint64 `gorm:"primary_key;autoIncrement; not null;" json:"id"`
	OrgID       uint64 `json:"orgID"`
	CreatorID   uint64 `gorm:"not null" json:"creatorID"`
	CreatorName string `gorm:"not null" json:"creatorName"`
	Name        string `gorm:"size:32; not null;" json:"name"`
	Profile     string `gorm:"not null;" json:"profile"`
	Version     int    `gorm:"not null;" json:"version"`
	Readable    int    `gorm:"not null" json:"readable"`
	Writable    int    `gorm:"not null" json:"writable"`
}

// 关系表

// Invitation 用户组织关系
type Invitation struct {
	ID        uint64 `gorm:"primary_key; autoIncrement;not null;" json:"id"`
	UserID    uint64 `gorm:"not null;" json:"userID"`
	UserName  string `gorm:"not null;" json:"userName"`
	UserEmail string `gorm:"not null;" json:"userEmail"`
	OrgID     uint64 `gorm:"not null;" json:"orgID"`
	OrgName   string `gorm:"not null;" json:"orgName"`
	IsAdmin   bool   `gorm:"not null;" json:"isAdmin"`
	IsValid   bool   `gorm:"not null;default:false" json:"isValid"`
}

// ContestProblem 比赛问题关系
type ContestProblem struct {
	ID          uint64 `gorm:"primary_key; autoIncrement;not null;" json:"id"`
	ContestID   uint64 `gorm:"not null;" json:"contestID"`
	ProblemID   uint64 `gorm:"not null;" json:"problemID"`
	ProblemName string `gorm:"not null;" json:"problemName"`
	Difficulty  int    `gorm:"not null;" json:"difficulty"`
}

// Result 用户问题关系
type Result struct {
	ID          uint64    `gorm:"primary_key; autoIncrement;not null;" json:"id"`
	UserID      uint64    `gorm:"not null;" json:"userID"`
	ProblemID   uint64    `gorm:"not null;" json:"problemID"`
	Result      int       `gorm:"not null;" json:"result"` // 0 AC , 1 WA , 2 TLE, 3 RE
	CreatedTime time.Time `gorm:"autoCreateTime;" json:"createdTime"`
}
