package repository

import (
	"errors"
	"fmt"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

type InstallRepository interface {
	CreateTmpDB(setting *domain.Setting) (*gorm.DB, error)
	CreateDB() error
}

type installRepository struct {
}

func (repo *installRepository) CreateTmpDB(setting *domain.Setting) (*gorm.DB, error) {
	return repo.connectDB(setting)
}

func (repo *installRepository) CreateDB() error {
	//TODO implement me
	panic("implement me")
}

func (repo *installRepository) connectDB(setting *domain.Setting) (*gorm.DB, error) {
	engine := strings.ToLower(setting.DB.Engine)
	switch engine {
	case "mysql":
		return gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
			setting.DB.User,
			setting.DB.Password,
			setting.DB.Host,
			setting.DB.Port,
			setting.DB.Database,
			setting.DB.Charset,
		)))
	case "postgres":
		return gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d",
			setting.DB.Host,
			setting.DB.User,
			setting.DB.Database,
			setting.DB.Password,
			setting.DB.Port,
		)))
	}
	return nil, errors.New("engine error")
}
