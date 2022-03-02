package database

// UserOrgRel 用户组织关系
type UserOrgRel struct {
	ID     uint64 `gorm:"primary_key; autoIncrement;not null;" json:"id"`
	UserID uint64 `gorm:"not null;" json:"user_id"`
	OrgID  uint64 `gorm:"not null;" json:"org_id"`
}

// CompOrgRel 比赛组织关系
type CompOrgRel struct {
	ID     uint64 `gorm:"primary_key; autoIncrement;not null;" json:"id"`
	CompID uint64 `gorm:"not null;" json:"comp_id"`
	OrgID  uint64 `gorm:"not null;" json:"org_id"`
}

// CompProbRel 比赛问题关系
type CompProbRel struct {
	ID     uint64 `gorm:"primary_key; autoIncrement;not null;" json:"id"`
	CompID uint64 `gorm:"not null;" json:"comp_id"`
	ProbID uint64 `gorm:"not null;" json:"prob_id"`
}
