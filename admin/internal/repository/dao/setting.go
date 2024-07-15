package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type SettingDAO interface {
	FindByKey(ctx context.Context, key string) ([]Setting, error)
	Insert(ctx context.Context, set Setting) (id int64, err error)
	UpDate(ctx context.Context, set Setting) error
	Delete(ctx context.Context, id int64) error
}

type GORMSettingDAO struct {
	db *gorm.DB
}

func NewGORMSettingDAO(db *gorm.DB) SettingDAO {
	return &GORMSettingDAO{
		db: db,
	}
}

func (dao *GORMSettingDAO) FindByKey(ctx context.Context, key string) ([]Setting, error) {
	var settings []Setting
	err := dao.db.WithContext(ctx).
		Where("key = ?", key).
		Find(&settings).Error
	return settings, err
}

func (dao *GORMSettingDAO) Insert(ctx context.Context, set Setting) (int64, error) {
	now := time.Now().UnixMilli()
	set.Ctime = now
	set.Utime = now
	err := dao.db.WithContext(ctx).Create(&set).Error
	return set.Id, err
}

func (dao *GORMSettingDAO) UpDate(ctx context.Context, set Setting) error {
	set.Utime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Updates(&set).Error
}

func (dao *GORMSettingDAO) Delete(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Where("id = ?", id).Delete(&Setting{}).Error
}

type Setting struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Code  string `gorm:"type:varchar(32)"`
	Key   string `gorm:"type:varchar(64)"`
	Value string `gorm:"type:varchar(4096)"`

	Utime int64
	Ctime int64
}
