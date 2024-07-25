package dao

import (
	"context"
	"gorm.io/gorm"
)

type SettingDAO interface {
	FindByKey(ctx context.Context, code string) ([]Setting, error)
}

type GORMSettingDAO struct {
	db *gorm.DB
}

func NewSettingDAO(db *gorm.DB) SettingDAO {
	return &GORMSettingDAO{
		db: db,
	}
}

func (dao *GORMSettingDAO) FindByKey(ctx context.Context, code string) ([]Setting, error) {
	var settings []Setting
	err := dao.db.WithContext(ctx).Where("`code` = ?", code).Find(&settings).Error
	return settings, err
}

type Setting struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Code  string `gorm:"type:varchar(32)"`
	Key   string `gorm:"type:varchar(64)"`
	Value string `gorm:"type:varchar(4096)"`

	Utime int64
	Ctime int64
}
