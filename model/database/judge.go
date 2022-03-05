package database

import "time"

// Competition 比赛
type Competition struct {
	ID        uint64    `gorm:"primary_key; autoIncrement; not null;" json:"id"`
	OrgID     uint64    `gorm:"not null;" json:"orgID"` // 比赛所属的组织ID
	Name      string    `gorm:"size:32; not null;" json:"name"`
	Profile   string    `json:"profile"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// Problem 题目
type Problem struct {
	ID           uint64    `gorm:"primary_key;autoIncrement;not null;" json:"id"`
	Name         string    `gorm:"size:32; not null" json:"name"`
	Version      int       `gorm:"not null;" json:"version"`
	Difficulty   int       `gorm:"not null" json:"difficulty"`
	Readable     int       `gorm:"not null" json:"readable"`
	Writable     int       `gorm:"not null" json:"writable"`
	Organization uint64    `json:"organization"`
	Creator      uint64    `gorm:"not null" json:"creator"`
	CreatedTime  time.Time `gorm:"autoCreateTime" json:"createdTime"`
}
