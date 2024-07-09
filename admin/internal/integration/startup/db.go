package startup

import (
	"github.com/zc-zht/super-job/admin/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitTestDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root123@tcp(localhost:3306)/super_job"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
