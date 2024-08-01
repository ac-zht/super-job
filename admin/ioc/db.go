package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root123@tcp(localhost:3306)/super_job"))
	if err != nil {
		panic(err)
	}
	//err = dao.InitTables(db)
	//if err != nil {
	//	panic(err)
	//}
	return db
}
