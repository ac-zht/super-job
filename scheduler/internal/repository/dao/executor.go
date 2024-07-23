package dao

type Executor struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Name  string `gorm:"type:varchar(256);unique"`
	Hosts string `gorm:"type:varchar(512)"`
}
