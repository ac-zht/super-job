package dao

import (
	"context"
	"time"
)

type LoginLogDAO interface {
	Insert(ctx context.Context, l LoginLog) (int64, error)
	List(ctx context.Context, offset, limit int) ([]LoginLog, error)
	Total(ctx context.Context) (int64, error)
}

type GORMLoginLogDAO struct {
	BaseModel
}

func NewGORMLoginLogDAO(base BaseModel) LoginLogDAO {
	return &GORMLoginLogDAO{BaseModel: base}
}

func (dao *GORMLoginLogDAO) Insert(ctx context.Context, l LoginLog) (int64, error) {
	l.Ctime = time.Now().UnixMilli()
	err := dao.DB().WithContext(ctx).Create(&l).Error
	return l.Id, err
}

func (dao *GORMLoginLogDAO) List(ctx context.Context, offset, limit int) ([]LoginLog, error) {
	var logs []LoginLog
	err := dao.DB().WithContext(ctx).Offset(offset).Limit(limit).Find(&logs).Error
	return logs, err
}

func (dao *GORMLoginLogDAO) Total(ctx context.Context) (int64, error) {
	var total int64
	err := dao.DB().WithContext(ctx).Model(&LoginLog{}).Count(&total).Error
	return total, err
}

type LoginLog struct {
	Id       int64
	Username string
	Ip       string
	Ctime    int64
}
