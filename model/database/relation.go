package database

// UserOrgRel 用户组织关系
type UserOrgRel struct {
	ID      uint64 `gorm:"primary_key; autoIncrement;not null;" json:"id"`
	UserID  uint64 `gorm:"not null;" json:"userID"`
	OrgID   uint64 `gorm:"not null;" json:"orgID"`
	IsAdmin bool   `gorm:"not null;" json:"isAdmin"`
}

// CompProbRel 比赛问题关系
type CompProbRel struct {
	ID     uint64 `gorm:"primary_key; autoIncrement;not null;" json:"id"`
	CompID uint64 `gorm:"not null;" json:"compID"`
	ProbID uint64 `gorm:"not null;" json:"probID"`
}
