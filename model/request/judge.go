package request

// 比赛
type Competition struct {
	ID   uint64 `gorm:"primary_key; not null;" json:"id"`
	Name string `gorm:"size:32; not null; unique" json:"name"`
}

// 题目
type Problem struct {
	ID   uint64 `gorm:"primary_key; not null;" json:"id"`
	Name string `gorm:"size:32; not null; unique" json:"name"`
}
