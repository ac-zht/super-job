package dao

type InstallDAO interface {
	InitTables() error
}

type GORMInstallDAO struct {
	BaseModel
}

func NewInstallDAO(base BaseModel) InstallDAO {
	return &GORMInstallDAO{
		BaseModel: base,
	}
}

func (dao *GORMInstallDAO) InitTables() error {
	return dao.DB().AutoMigrate(
		&User{},
		&Task{},
		&TaskLog{},
		&Executor{},
		&Setting{},
	)
}
