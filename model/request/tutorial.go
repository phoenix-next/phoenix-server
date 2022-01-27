package request

// 教程
type Tutorial struct {
	ID   uint64 `gorm:"primary_key; not null;" json:"id"`
	Name string `gorm:"size:32; not null; unique" json:"name"`
}
