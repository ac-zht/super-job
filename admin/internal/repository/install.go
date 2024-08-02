package repository

import (
	"errors"
	"fmt"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository/dao"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
)

type app struct {
	//ConfDir 配置文件目录
	ConfDir string
	//Config 应用配置文件
	Config string
	//Installed 应用是否已安装
	Installed bool
	//Setting 应用配置
	Setting *domain.Setting
	Mode    string
}

var App = app{Setting: &domain.Setting{}}

type InstallRepository interface {
	PingDB(setting *domain.Setting) error
	CreateDB() *gorm.DB
	InitTables() error
}

type installRepository struct {
	dao dao.InstallDAO
}

func NewInstallRepository(dao dao.InstallDAO) InstallRepository {
	return &installRepository{
		dao: dao,
	}
}

func (repo *installRepository) PingDB(setting *domain.Setting) error {
	db, err := repo.connectDB(setting)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	err = sqlDB.Ping()
	if err != nil {
		switch setting.DB.Engine {
		case "mysql":
			mysqlError, ok := err.(*mysqlDriver.MySQLError)
			if ok && mysqlError.Number == 1049 {
				err = errors.New("database not exist")
			}
			return err
		case "postgres":
			pgError, ok := err.(*pq.Error)
			if ok && pgError.Code == "3D000" {
				err = errors.New("database not exist")
			}
			return err
		}
	}
	return err
}

func (repo *installRepository) CreateDB() *gorm.DB {
	db, err := repo.connectDB(App.Setting)
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("创建gorm引擎失败#%v", err))
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(App.Setting.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(App.Setting.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(dao.DbMaxLifeTime)
	return db
}

func (repo *installRepository) connectDB(setting *domain.Setting) (*gorm.DB, error) {
	engine := strings.ToLower(setting.DB.Engine)
	config := &gorm.Config{}
	if App.Setting.DB.Prefix != "" {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix: fmt.Sprintf("%s_", setting.DB.Prefix),
		}
	}
	config.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		LogLevel: logger.Info,
	})
	switch engine {
	case "mysql":
		return gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
			setting.DB.User,
			setting.DB.Password,
			setting.DB.Host,
			setting.DB.Port,
			setting.DB.Database,
			setting.DB.Charset,
		)), config)
	case "postgres":
		return gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d",
			setting.DB.Host,
			setting.DB.User,
			setting.DB.Database,
			setting.DB.Password,
			setting.DB.Port,
		)), config)
	}
	return nil, errors.New("engine error")
}

func (repo *installRepository) InitTables() error {
	return repo.dao.InitTables()
}
