package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var ErrNoMoreJob = gorm.ErrRecordNotFound

type JobDAO interface {
	List(ctx context.Context, offset, limit int) ([]Job, error)
	Insert(ctx context.Context, j Job) error
	Delete(ctx context.Context, id int) error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (dao *GORMJobDAO) List(ctx context.Context, offset, limit int) ([]Job, error) {
	var jobs []Job
	err := dao.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&jobs).Error
	return jobs, err
}

func (dao *GORMJobDAO) Delete(ctx context.Context, id int) error {
	return dao.db.WithContext(ctx).Where("id = ?", id).Delete(&Job{}).Error
}

func (dao *GORMJobDAO) Insert(ctx context.Context, j Job) error {
	now := time.Now().UnixMilli()
	j.Ctime = now
	j.Utime = now
	return dao.db.WithContext(ctx).Create(&j).Error
}

type Job struct {
	Id         int64 `gorm:"primaryKey,autoIncrement"`
	ExecId     int64
	Name       string `gorm:"type:varchar(256);unique"`
	Protocol   uint8  `json:"protocol" gorm:"tinyint"` // 协议 1:http 2:rpc 3:系统命令
	Cfg        string
	Expression string
	Version    int64
	//可建next_time和status的联合索引
	NextTime int64 `gorm:"index:status_next_index"`
	Status   int   `gorm:"index:status_next_index"`
	Ctime    int64
	Utime    int64
}

const (
	jobStatusWaiting = iota
	jobStatusRunning
)
