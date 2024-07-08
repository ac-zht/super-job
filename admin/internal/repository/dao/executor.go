package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ExecutorDAO interface {
	List(ctx context.Context, offset, limit int) ([]Executor, error)
	Insert(ctx context.Context, exec Executor) (int64, error)
	Update(ctx context.Context, exec Executor) error
	Delete(ctx context.Context, id int64) error
}

type GORMExecutorDAO struct {
	db *gorm.DB
}

func NewExecutorDAO(db *gorm.DB) ExecutorDAO {
	return &GORMExecutorDAO{
		db: db,
	}
}

func (dao GORMExecutorDAO) List(ctx context.Context, offset, limit int) ([]Executor, error) {
	var execs []Executor
	err := dao.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&execs).Error
	return execs, err
}

func (dao GORMExecutorDAO) Insert(ctx context.Context, exec Executor) (int64, error) {
	now := time.Now().UnixMilli()
	exec.Ctime = now
	exec.Utime = now
	err := dao.db.WithContext(ctx).Create(&exec).Error
	return exec.Id, err
}

func (dao GORMExecutorDAO) Update(ctx context.Context, exec Executor) error {
	return dao.db.WithContext(ctx).Updates(&exec).Error
}

func (dao GORMExecutorDAO) Delete(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Where("id = ?", id).Delete(&Executor{}).Error
}

type Executor struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Name  string `gorm:"type:varchar(256);unique"`
	Hosts string `gorm:"type:varchar(512)"`
	Ctime int64
	Utime int64
}
