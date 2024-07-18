package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var ErrNoMoreTask = gorm.ErrRecordNotFound

type TaskDAO interface {
	List(ctx context.Context, offset, limit int) ([]Task, error)
	GetById(ctx context.Context, id int64) (Task, error)
	Insert(ctx context.Context, j Task) (int64, error)
	Update(ctx context.Context, j Task) error
	Delete(ctx context.Context, id int64) error
}

type GORMTaskDAO struct {
	db *gorm.DB
}

func NewTaskDAO(db *gorm.DB) TaskDAO {
	return &GORMTaskDAO{
		db: db,
	}
}

func (dao *GORMTaskDAO) List(ctx context.Context, offset, limit int) ([]Task, error) {
	var tasks []Task
	err := dao.db.WithContext(ctx).
		Preload("Executor").
		Offset(offset).
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

func (dao *GORMTaskDAO) GetById(ctx context.Context, id int64) (Task, error) {
	var task Task
	err := dao.db.WithContext(ctx).First(&task, id).Error
	return task, err
}

func (dao *GORMTaskDAO) Delete(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Where("id = ?", id).Delete(&Task{}).Error
}

func (dao *GORMTaskDAO) Insert(ctx context.Context, j Task) (int64, error) {
	now := time.Now().UnixMilli()
	j.Ctime = now
	j.Utime = now
	err := dao.db.WithContext(ctx).Create(&j).Error
	return j.Id, err
}

func (dao *GORMTaskDAO) Update(ctx context.Context, j Task) error {
	j.Utime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Updates(&j).Error
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
	ExecutorHandler string
	Command         string

	//失败重试策略
	Timeout       int64
	RetryTimes    int64
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
