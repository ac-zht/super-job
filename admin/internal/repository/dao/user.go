package dao

type User struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	Name    string `gorm:"type:varchar(256);unique"`
	Email   string `gorm:"type:varchar(256);unique"`
	Salt    string
	IsAdmin uint8
	Status  uint8
	Ctime   int64
	Utime   int64
}
