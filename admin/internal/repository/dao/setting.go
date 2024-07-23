package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type SettingDAO interface {
	FindByKey(ctx context.Context, code string) ([]Setting, error)
	Insert(ctx context.Context, setting Setting) (int64, error)
	UpdateByCodeKey(ctx context.Context, code, key, value string) error
	Update(ctx context.Context, setting Setting) (int64, error)
	Delete(ctx context.Context, id int64) error
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

func (dao *GORMSettingDAO) Insert(ctx context.Context, setting Setting) (int64, error) {
	now := time.Now().UnixMilli()
	setting.Ctime = now
	setting.Utime = now
	err := dao.db.WithContext(ctx).Create(&setting).Error
	return setting.Id, err
}

func (dao *GORMSettingDAO) Update(ctx context.Context, setting Setting) (int64, error) {
	setting.Utime = time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Updates(&setting).Error
	return setting.Id, err
}

func (dao *GORMSettingDAO) UpdateByCodeKey(ctx context.Context, code, key, value string) error {
	return dao.db.WithContext(ctx).Model(&Setting{}).Where("`code` = ? AND `key` = ?", code, key).Update("value", value).Error
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
