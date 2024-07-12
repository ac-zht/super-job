package dao

type Setting struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Code  string `gorm:"type:varchar(32)"`
	Key   string `gorm:"type:varchar(64)"`
	Value string `gorm:"type:varchar(4096)"`
}
