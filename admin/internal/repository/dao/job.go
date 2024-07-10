package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var ErrNoMoreJob = gorm.ErrRecordNotFound

type JobDAO interface {
	List(ctx context.Context, offset, limit int) ([]Job, error)
	Insert(ctx context.Context, j Job) (int64, error)
	Update(ctx context.Context, j Job) error
	Delete(ctx context.Context, id int64) error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func NewJobDAO(db *gorm.DB) JobDAO {
	return &GORMJobDAO{
		db: db,
	}
}

func (dao *GORMJobDAO) List(ctx context.Context, offset, limit int) ([]Job, error) {
	var jobs []Job
	err := dao.db.WithContext(ctx).
		Preload("Executor").
		Offset(offset).
		Limit(limit).
		Find(&jobs).Error
	return jobs, err
}

func (dao *GORMJobDAO) Delete(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Where("id = ?", id).Delete(&Job{}).Error
}

func (dao *GORMJobDAO) Insert(ctx context.Context, j Job) (int64, error) {
	now := time.Now().UnixMilli()
	j.Ctime = now
	j.Utime = now
	err := dao.db.WithContext(ctx).Create(&j).Error
	return j.Id, err
}

func (dao *GORMJobDAO) Update(ctx context.Context, j Job) error {
	j.Utime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Updates(&j).Error
}

type Job struct {
	Id         int64 `gorm:"primaryKey,autoIncrement"`
	ExecId     int64
	Name       string `gorm:"type:varchar(256);unique"`
	Cfg        string
	Expression string
	Status     uint8 `gorm:"index:status_next_index"`
	//可建next_time和status的联合索引
	NextTime   int64 `gorm:"index:status_next_index"`
	Protocol   uint8 `gorm:"tinyint"` // 协议 1:http 2:rpc 3:系统命令
	HttpMethod uint8 `gorm:"tinyint"`
	Multi      uint8 //该任务同一时间是否只运行在一个实例上

	//失败重试策略
	Timeout       int64
	RetryTimes    int64
	RetryInterval int64

	//消息通知
	NotifyStatus     uint8 `gorm:"tinyint"`
	NotifyType       uint8 `gorm:"tinyint"`
	NotifyReceiverId int64
	NotifyKeyword    string

	Version int64

	Creator int64
	//最后一次更新的人
	Updater int64

	Ctime    int64
	Utime    int64
	Executor Executor `gorm:"foreignKey:ExecId"`
}

const (
	jobStatusWaiting = iota
	jobStatusRunning
)
