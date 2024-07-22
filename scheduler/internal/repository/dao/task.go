package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var ErrNoMoreTask = gorm.ErrRecordNotFound

type TaskDAO interface {
	Preempt(ctx context.Context) (Task, error)
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
	Insert(ctx context.Context, j Task) error
}

type GORMTaskDAO struct {
	db *gorm.DB
}

func (dao *GORMTaskDAO) Preempt(ctx context.Context) (Task, error) {
	db := dao.db.WithContext(ctx)
	for {
		now := time.Now().UnixMilli()
		var j Task
		err := db.Where(
			"next_time <= ? AND status = ?",
			now, taskStatusWaiting).First(&j).Error
		if err != nil {
			return Task{}, err
		}
		res := db.Model(&Task{}).
			Where("id = ? AND version=?", j.Id, j.Version).
			Updates(map[string]any{
				"utime":   now,
				"version": j.Version + 1,
				"status":  taskStatusRunning,
			})
		if res.Error != nil {
			return Task{}, res.Error
		}
		if res.RowsAffected == 1 {
			return j, nil
		}
	}
}

func (dao *GORMTaskDAO) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	return dao.db.WithContext(ctx).Model(&Task{}).
		Where("id=?", id).Updates(map[string]any{
		"utime":     time.Now().UnixMilli(),
		"next_time": t.UnixMilli(),
	}).Error
}

func (dao *GORMTaskDAO) UpdateUtime(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&Task{}).
		Where("id=?", id).Updates(map[string]any{
		"utime": time.Now().UnixMilli(),
	}).Error
}

func (dao *GORMTaskDAO) Release(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&Task{}).
		Where("id = ?", id).Updates(map[string]any{
		"status": taskStatusWaiting,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

type Task struct {
	Id         int64 `gorm:"primaryKey,autoIncrement"`
	ExecId     int64
	Name       string `gorm:"type:varchar(256);unique"`
	Cfg        string
	Expression string `gorm:"type:varchar(256)"`
	//可建next_time和status的联合索引
	NextTime   int64 `gorm:"index:next_status_index"`
	Status     uint8 `gorm:"index:next_status_index"`
	Protocol   uint8 `gorm:"tinyint"` // 协议 1:http 2:rpc 3:系统命令
	HttpMethod uint8 `gorm:"tinyint"`
	Multi      uint8 //该任务同一时间是否只运行在一个实例上
	//方法或命令
	ExecutorHandler       string
	Command               string
	ExecutorRouteStrategy string `gorm:"type:varchar(50)"`

	//失败重试策略
	Timeout       int64
	RetryTimes    int8
	RetryInterval int64

	//消息通知
	NotifyStatus     uint8  `gorm:"tinyint"`           // 任务执行结束是否通知 0: 不通知 1: 失败通知 2: 执行结束通知 3: 任务执行结果关键字匹配通知
	NotifyType       uint8  `gorm:"tinyint"`           // 通知类型 1: 邮件 2: slack 3: webhook
	NotifyReceiverId string `gorm:"type:varchar(256)"` // 通知接受者ID, setting表主键ID，多个ID逗号分隔
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
	taskStatusWaiting = iota
	taskStatusRunning
)
