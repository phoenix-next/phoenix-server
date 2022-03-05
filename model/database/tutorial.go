package database

// Tutorial 教程
type Tutorial struct {
	ID           uint64 `gorm:"primary_key;autoIncrement; not null;" json:"id"`
	Name         string `gorm:"size:32; not null;" json:"name"`
	Profile      string `gorm:"not null;" json:"profile"`
	Version      int    `gorm:"not null;" json:"version"`
	Readable     int    `gorm:"not null" json:"readable"`
	Writable     int    `gorm:"not null" json:"writable"`
	Organization uint64 `json:"organization"`
	Creator      uint64 `gorm:"not null" json:"creator"`
}
