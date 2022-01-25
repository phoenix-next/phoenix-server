package model

// 用户
type User struct {
	ID       uint64 `gorm:"primary_key; not null;" json:"id"`
	Name     string `gorm:"size:32; not null; unique" json:"name"`
	Email    string `gorm:"size:32; not null; unique;" json:"email"`
	Password string `gorm:"size:32; not null;" json:"password"`
	Avatar   string `gorm:"size:256;" json:"avatar"`
	Profile  string `gorm:"size:255;" json:"profile"`
}
