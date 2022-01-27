package model

import "time"

// 用户
type User struct {
	ID       uint64    `gorm:"primary_key; not null;" json:"id"`
	Name     string    `gorm:"size:32; not null; unique" json:"name"`
	Email    string    `gorm:"size:32; not null; unique;" json:"email"`
	Password string    `gorm:"size:32; not null;" json:"password"`
	Avatar   string    `gorm:"size:256;" json:"avatar"`
	Profile  string    `gorm:"size:255;" json:"profile"`
	RegTime  time.Time `gorm:"datetime" json:"reg_time"`
	UserType uint64    `gorm:"default:0;" json:"user_type"` // 0为普通用户，1为至少为某一组织的管理员
}

// 组织
type Organization struct {
	ID          uint64    `gorm:"primary_key; not null;" json:"id"`
	Name        string    `gorm:"size:32; not null; unique" json:"name"`
	Profile     string    `gorm:"size:255;" json:"profile"`
	CreatorID   uint64    `gorm:"not null;" json:"creator_id"` // 创建者ID
	CreatorName string    `gorm:"size:32;not null;" json:"creator_name"`
	CreatedTime time.Time `gorm:"datetime;" json:"created_time"`
}
