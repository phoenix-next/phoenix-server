package request

// 帖子
type Post struct {
	ID   uint64 `gorm:"primary_key; not null;" json:"id"`
	Name string `gorm:"size:32; not null; unique" json:"name"`
}

// 评论
type Comment struct {
	ID   uint64 `gorm:"primary_key; not null;" json:"id"`
	Name string `gorm:"size:32; not null; unique" json:"name"`
}
