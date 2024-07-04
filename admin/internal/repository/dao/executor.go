package dao

import (
	"context"
	"gorm.io/gorm"
)

type ExecutorDAO interface {
	List(ctx context.Context, offset, limit int) ([]Executor, error)
	Insert(ctx context.Context, exec Executor) error
	Update(ctx context.Context, exec Executor) error
	Delete(ctx context.Context, id int64) error
}

type GORMExecutorDAO struct {
	db *gorm.DB
}

func (G GORMExecutorDAO) List(ctx context.Context, offset, limit int) ([]Executor, error) {
	//TODO implement me
	panic("implement me")
}

func (G GORMExecutorDAO) Insert(ctx context.Context, exec Executor) error {
	//TODO implement me
	panic("implement me")
}

func (G GORMExecutorDAO) Update(ctx context.Context, exec Executor) error {
	//TODO implement me
	panic("implement me")
}

func (G GORMExecutorDAO) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

type Executor struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Name  string `gorm:"type:varchar(256);unique"`
	Hosts string `gorm:"type:varchar(512)"`
}
