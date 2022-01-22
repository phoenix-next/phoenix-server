package model

import (
	"time"
)

// 用户
type User struct {
	UserID       uint64    `gorm:"primary_key; not null;" json:"user_id"`
	Username     string    `gorm:"size:32; not null; unique" json:"username"`
	Password     string    `gorm:"size:32; not null" json:"password"`
	UserInfo     string    `gorm:"size:255;" json:"user_info"`
	UserType     uint64    `gorm:"default:0" json:"user_type"` // 0: 普通用户，1: 认证机构用户,2 管理员
	Avatar       string    `gorm:"size:256" json:"avatar"`
	Email        string    `gorm:"size:32;" json:"email"`
	HasConfirmed bool      `gorm:"default:false" json:"has_confirmed"`
	RegTime      time.Time `gorm:"column:reg_time;type:datetime" json:"reg_time"`
}
