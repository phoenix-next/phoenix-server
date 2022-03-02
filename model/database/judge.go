package database

import "time"

// Competition 比赛
type Competition struct {
	ID    uint64 `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	OrgID uint64 `gorm:"not null;" json:"org_id"` // 所属比赛ID
	Name  string `gorm:"size:32; not null; unique" json:"name"`
}

// Problem 题目
type Problem struct {
	ID           uint64    `gorm:"primary_key; not null;" json:"id"`
	Name         string    `gorm:"size:32; not null; unique" json:"name"`
	Version      int       `gorm:"autoIncrement;not null;" json:"version"`
	Difficulty   int       `gorm:"not null" json:"difficulty"`
	Readable     int       `gorm:"not null" json:"readable"`
	Writable     int       `gorm:"not null" json:"writable"`
	Organization uint64    `gorm:"not null" json:"organization"`
	Creator      uint64    `gorm:"not null" json:"creator"`
	CreatedTime  time.Time `gorm:"autoCreateTime" json:"created_time"`
}
