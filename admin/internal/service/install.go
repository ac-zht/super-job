package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository"
	"github.com/ac-zht/super-job/admin/internal/repository/dao"
	"github.com/ac-zht/super-job/admin/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strconv"
)

const (
	DEV  = "development"
	PROD = "production"
	TEST = "test"
)

var App = &repository.App

func SetAppSetting(setting *domain.Setting) {
	App.Setting = setting
}

type InstallService interface {
	CreateDB() *gorm.DB
	Store(ctx context.Context, ins domain.Installation) error
	Status(ctx context.Context) (bool, error)
	CreateInstallLock() error
}

type InstallSvc struct {
	settingRepo    repository.SettingRepository
	installRepo    repository.InstallRepository
	webSettingRepo repository.WebSettingRepository
	userRepo       repository.UserRepository
}

func NewInstallService(setRepo repository.SettingRepository,
	installRepo repository.InstallRepository,
	webSettingRepo repository.WebSettingRepository,
	userRepo repository.UserRepository) InstallService {
	return &InstallSvc{
		settingRepo:    setRepo,
		installRepo:    installRepo,
		webSettingRepo: webSettingRepo,
		userRepo:       userRepo,
	}
}

func (svc *InstallSvc) CreateDB() *gorm.DB {
	return svc.installRepo.CreateDB()
}

func (svc *InstallSvc) Store(ctx context.Context, install domain.Installation) error {
	//ping数据库
	err := svc.pingDB(ctx, install)
	if err != nil {
		return errors.New(fmt.Sprintf("ping database fail#%v", err))
	}
	//根据提交的消息写数据库配置文件
	err = svc.writeConfig(ctx, install)
	if err != nil {
		return errors.New(fmt.Sprintf("config wirte to file fail#%v", err))
	}
	//读取文件到内存
	appConfig, err := svc.webSettingRepo.Read(App.Config)
	if err != nil {
		return errors.New(fmt.Sprintf("read app config fail#%v", err))
	}
	App.Setting = appConfig
	//创建全局db连接
	dao.SetGlobalDB(svc.CreateDB())
	//创建数据表
	err = svc.installRepo.InitTables()
	if err != nil {
		return errors.New(fmt.Sprintf("create table fail#%v", err))
	}
	//初始化配置表字段
	err = svc.settingRepo.InitBasicField(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("init basic field fail#%v", err))
	}
	//创建管理员
	err = svc.creatAdminUser(ctx, install)
	if err != nil {
		return errors.New(fmt.Sprintf("create administrator fail#%v", err))
	}
	err = svc.CreateInstallLock()
	if err != nil {
		return errors.New(fmt.Sprintf("create install lock fail#%v", err))
	}
	App.Installed = true
	return nil
}

func (svc *InstallSvc) pingDB(ctx context.Context, ins domain.Installation) error {
	var s domain.Setting
	s.DB.Engine = ins.DbType
	s.DB.Host = ins.DbHost
	s.DB.Port = ins.DbPort
	s.DB.User = ins.DbUsername
	s.DB.Password = ins.DbPassword
	s.DB.Database = ins.DbName
	s.DB.Charset = "utf8"
	return svc.installRepo.PingDB(&s)
}

func (svc *InstallSvc) writeConfig(ctx context.Context, ins domain.Installation) error {
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
		//"concurrency.queue", "500",
		"auth_secret", utils.RandAuthToken(),
		"ca_file", "",
		"cert_file", "",
		"key_file", "",
	}
	return svc.webSettingRepo.Write(dbConfig, App.Config)
}

func (svc *InstallSvc) Status(ctx context.Context) (bool, error) {
	return App.Installed, nil
}

func (svc *InstallSvc) creatAdminUser(ctx context.Context, install domain.Installation) error {
	user := domain.User{
		Name:     install.AdminUsername,
		Email:    install.AdminEmail,
		Password: install.AdminPassword,
		IsAdmin:  2,
	}
	_, err := svc.userRepo.Create(ctx, user)
	return err
}

func (svc *InstallSvc) CreateInstallLock() error {
	_, err := os.Create(filepath.Join(App.ConfDir, "/install.lock"))
	if err != nil {
		zap.L().Error("创建安装锁文件conf/install.lock失败")
	}
	return err
}
