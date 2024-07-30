package service

import (
	"context"
	"errors"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository"
	"github.com/ac-zht/super-job/admin/pkg/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"strconv"
)

type InstallService interface {
	Store(ctx context.Context, ins domain.Installation) error
	Status(ctx context.Context) (bool, error)
}

type installService struct {
	settingRepo    repository.SettingRepository
	installRepo    repository.InstallRepository
	webSettingRepo repository.WebSettingRepository
}

func NewInstallService(setRepo repository.SettingRepository, installRepo repository.InstallRepository, webSettingRepo repository.WebSettingRepository) InstallService {
	return &installService{
		settingRepo:    setRepo,
		installRepo:    installRepo,
		webSettingRepo: webSettingRepo,
	}
}

func (svc *installService) Store(ctx context.Context, install domain.Installation) error {
	//ping数据库
	err := svc.pingDB(ctx, install)
	if err != nil {
		return err
	}
	//根据提交的消息写数据库配置文件
	err = svc.writeConfig(ctx, install)
	if err != nil {
		return err
	}
	//读取文件到内存
	//创建数据库
	//生成表
	//初始化配置表字段
	err = svc.settingRepo.InitBasicField(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (svc *installService) pingDB(ctx context.Context, ins domain.Installation) error {
	var s domain.Setting
	s.DB.Engine = ins.DbType
	s.DB.Host = ins.DbHost
	s.DB.Port = ins.DbPort
	s.DB.User = ins.DbUsername
	s.DB.Password = ins.DbPassword
	s.DB.Database = ins.DbName
	s.DB.Charset = "utf8"
	db, err := svc.installRepo.CreateTmpDB(&s)
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
		switch s.DB.Engine {
		case "mysql":
			mysqlError, ok := err.(*mysql.MySQLError)
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

func (svc *installService) writeConfig(ctx context.Context, ins domain.Installation) error {
	dbConfig := []string{
		"db.engine", ins.DbType,
		"db.host", ins.DbHost,
		"db.port", strconv.Itoa(ins.DbPort),
		"db.user", ins.DbUsername,
		"db.password", ins.DbPassword,
		"db.database", ins.DbName,
		"db.prefix", ins.DbTablePrefix,
		"db.charset", "utf8",
		"db.max.idle.conns", "",
		"db.max.open.conns", "",
		"allow_ips", "",
		"app.name", "定时任务管理系统", // 应用名称
		"api.key", "",
		"api.secret", "",
		"enable_tls", "false",
		"concurrency.queue", "500",
		"auth_secret", utils.RandAuthToken(),
		"ca_file", "",
		"cert_file", "",
		"key_file", "",
	}
	return svc.webSettingRepo.Write(dbConfig, AppConfig)
}

func (svc *installService) Status(ctx context.Context) (bool, error) {
	return Installed, nil
}
