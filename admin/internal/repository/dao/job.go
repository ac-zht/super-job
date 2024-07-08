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
	Delete(ctx context.Context, id int) error
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

func (dao *GORMJobDAO) Delete(ctx context.Context, id int) error {
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
	return dao.db.WithContext(ctx).Updates(&j).Error
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
	Executor Executor `gorm:"foreignKey:ExecId"`
}

const (
	jobStatusWaiting = iota
	jobStatusRunning
)
