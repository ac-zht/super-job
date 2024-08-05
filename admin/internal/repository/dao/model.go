package dao

import (
	"github.com/ac-zht/gotools/option"
	"gorm.io/gorm"
	"time"
)

const (
	DbMaxLifeTime = time.Hour * 2
)

type CommonMap map[string]interface{}

var Db *gorm.DB

type BaseModel interface {
	DB() *gorm.DB
}

type BaseDbModel struct {
	Db *gorm.DB
}

func NewBaseModel(opts ...option.Option[BaseDbModel]) BaseModel {
	model := &BaseDbModel{}
	option.Apply[BaseDbModel](model, opts...)
	return model
}

func (m *BaseDbModel) DB() *gorm.DB {
	if m.Db != nil {
		return m.Db
	}
	return Db
}

func SetGlobalDB(db *gorm.DB) {
	Db = db
}
